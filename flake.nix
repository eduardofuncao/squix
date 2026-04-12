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
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "squix";
          version = "0.4.0-beta";

          src = ./.;

          # Run: nix build .#default 2>&1 | grep "got:" to get real hash
          vendorHash = "sha256-JRmNajvCb57dMo8eggOD1m4N01p2RSK8r49pmBB56Z0=";

          enableCGO = true;

          # Native dependencies (none needed - all drivers are pure Go)

          # Linker flags
          ldflags = [
            "-s"
            "-w"
            "-X main.Version=${self.packages.${system}.default.version}"
          ];

          meta = with pkgs.lib; {
            description = "Minimal CLI tool for managing SQL queries across multiple databases";
            homepage = "https://github.com/eduardofuncao/squix";
            license = licenses.mit;
            mainProgram = "squix";
          };
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
