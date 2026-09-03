{
  description = "Squix's SQL Stash - SQL query CLI tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        version = "0.5.4-beta"; # update when cutting a release

        mkSquix = { pname, tags ? [ ] }: pkgs.buildGoModule {
          inherit pname version tags;

          src = ./.;

          # Run: nix build .#default 2>&1 | grep "got:" to get real hash
          vendorHash = "sha256-kSv3VAQi+qdT29gZAjLmHauItaMFd9NG7bdRtQE1MZo=";

          # DuckDB is the only driver needing CGO; lite/minimal exclude it.
          enableCGO = ! (builtins.elem "noduckdb" tags);

          # Linker flags
          ldflags = [
            "-s"
            "-w"
            "-X main.Version=${version}"
          ];

          postInstall = pkgs.lib.optionalString (pname != "squix") ''
            mv "$out/bin/squix" "$out/bin/${pname}"
          '';

          meta = with pkgs.lib; {
            description = "Minimal CLI tool for managing SQL queries across multiple databases";
            homepage = "https://github.com/eduardofuncao/squix";
            license = licenses.mit;
            mainProgram = pname;
          };
        };
      in
      {
        packages.default = mkSquix { pname = "squix"; };

        # No DuckDB, Snowflake, Oracle
        packages.lite = mkSquix {
          pname = "squix-lite";
          tags = [ "noduckdb" "nosnowflake" "noracle" ];
        };

        # Only Postgres, MySQL, SQLite
        packages.minimal = mkSquix {
          pname = "squix-minimal";
          tags = [
            "noduckdb"
            "nosnowflake"
            "noracle"
            "noclickhouse"
            "nofirebird"
            "nosqlserver"
          ];
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            postgresql
          ];

          shellHook = ''
            echo "========================================="
            echo "Squix development environment ready!"
            echo "========================================="
            echo ""
            echo "Available tools:"
            echo "  - Go compiler"
            echo "  - PostgreSQL client (psql)"
            echo "  - SQLite client (sqlite3)"
            echo ""
          '';
        };
      }
    );
}
