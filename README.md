<div align="center">
<h1>
  <img width="auto" height="45" alt="image" src="https://github.com/user-attachments/assets/f82ceec8-9fc7-4253-9ec0-f5548c646996" />
  Squix's SQL Stash
<img width="auto" height="36" alt="image" src="https://github.com/user-attachments/assets/c128f28f-dd10-4213-9915-dedafe7ae831" />

</h1>
<img width="360" height="131" alt="image" src="https://github.com/user-attachments/assets/9428a75b-ffa4-4961-919b-e5ccf192ef26" />

### **SQL Query Stashing for Terminal Squirrels**

> **Bear Grylls:** "Out here in the wild database ecosystem, efficiency means survival. See that squirrel? That’s Squix, or _Sequillis termius_. He doesn’t panic-write queries under pressure. He prepares. He caches. He optimizes. While others are wrestling with joins in the dark, Squix already has his results. Extraordinary creature."

---

[![MIT License](https://img.shields.io/badge/license-MIT-white.svg)](LICENSE)
![go badge](https://img.shields.io/badge/g=Go-1.25+-00ADD8?%20logo=go&logoColor=white&label=go)
[![Matrix](https://img.shields.io/matrix/squix:matrix.org?server_fqdn=matrix.org&label=chat&color=green)](https://matrix.to/#/#squix-sql:matrix.org)

**A minimal CLI tool for managing and executing SQL queries across multiple databases. Written in Go, made beautiful with BubbleTea**

[Quick Start](#--------quick-start) • [Configuration](docs/configuration.md) • [Commands](docs/commands.md) • [Keybindings](docs/keybindings.md) • [Features](docs/features.md) • [Completion](docs/completion.md) • [Dbeesly](#-dbeesly) • [Roadmap](#--------roadmap) • [Contributing](#contributing)

> This project is currently in beta, please report unexpected behavior through the issues tab

</div>


---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle;" src="https://github.com/user-attachments/assets/464275ac-085e-451f-b783-c991d24d3635" />
    Demo
</h2>

![squixdemo2](https://github.com/user-attachments/assets/ee9653cf-6aaa-4be9-a898-37153ab0c898)

> Try out the [live demo](https://squix.live.eduardofuncao.com) (no install required!)

### Highlights

- **Query Library** - Save and organize your most-used queries
- **Runs in the CLI** - Execute queries with minimal overhead
- **Multi-Database** - Works with PostgreSQL, MySQL, SQLite, Oracle, SQL Server, ClickHouse and Firebird
- **Table view TUI** - Keyboard focused navigation with vim-style bindings
- **In-Place Editing** - Update cells, delete rows and edit your SQL directly from the results table
- **Export your data** - Export your data as CSV, JSON, SQL, Markdown or HTML tables

---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle;" src="https://github.com/user-attachments/assets/30765e98-13b3-4c18-81e7-faf224b60e0b" />
    Quick Start
</h2>

### Installation
Go to [the releases page](https://github.com/eduardofuncao/squix/releases) and find the correct version for your system. Download it and make sure the file is executable and moved to a directory in your $PATH.


<details>
<summary>Go install</summary>

Use go to install `squix` directly
```bash
go install github.com/eduardofuncao/squix/cmd/squix@latest
```
this will put the binary `squix` in your $GOBIN path (usually `~/go/bin`)
</details>

<details>
<summary>Build Manually</summary>

Follow these instructions to build the project locally
```bash
git clone https://github.com/eduardofuncao/squix

go build -o squix ./cmd/squix
```
The squix binary will be available in the root project directory
</details>

<details>
<summary>Nix / NixOS (Flake)</summary>

Squix is available as a Nix flake for easy installation on NixOS and systems with
Nix.


#### Run directly without installing
```bash
nix run github:eduardofuncao/squix
```

#### Install to user profile
```bash
nix profile install github:eduardofuncao/squix
```

#### Enter development shell
```bash
nix develop github:eduardofuncao/squix
```

#### NixOS System-wide

Add to your flake-based configuration.nix or flake.nix:

```nix
{
description = "My NixOS config";

inputs = {
  nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  squix.url = "github:eduardofuncao/squix";
};

outputs = { self, nixpkgs, squix, ... }: {
  nixosConfigurations.myHostname = nixpkgs.lib.nixosSystem {
    system = "x86_64-linux";
    modules = [
      {
        environment.systemPackages = [
          squix.packages.x86_64-linux.default
        ];
      }
    ];
  };
};
}
```

Then rebuild: sudo nixos-rebuild switch

#### Home Manager

Add to your home.nix or flake config:

```nix
{
inputs = {
  nixpkgs.url = "github:NixOS/nixpkgs/nix-unstable";
  squix.url = "github:eduardofuncao/squix";
};

outputs = { self, nixpkgs, squix, ... }: {
  homeConfigurations."username" = {
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
    modules = [
      {
        home.packages = [
          squix.packages.x86_64-linux.default
        ];
      }
    ];
  };
};
}
```

Then apply: home-manager switch
</details>

### Basic Usage

```bash
# Create your first connection (PostgreSQL example)
squix init mydb postgres "postgresql://user:pass@localhost:5432/mydb"

# Add a saved query
squix add list_users "SELECT * FROM users"

# List your saved queries
squix list queries

# Run it, this opens the interactive table viewer
squix run list_users

# Or run inline SQL
squix run "SELECT * FROM products WHERE price > 100"
```

### Navigating the Table

Once your query results appear, you can navigate and interact with the data:

```bash
# Use vim-style navigation or arrow-keys
j/k        # Move down/up
h/l        # Move left/right
g/G        # Jump to first/last row

# Copy data
y          # Yank (copy) current cell
v          # Enter visual mode to select multiple cells and copy with y
x          # Export selected data as csv, tsv, json, sql, markdown or html

# Edit data directly
u          # Update current cell (opens your $EDITOR)
D          # Delete current row

# Modify and re-run
e          # Edit the query and re-run it

# Exit
q          # Quit back to terminal
```

---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle" src="https://github.com/user-attachments/assets/8f5037c9-e616-4065-adfc-cd598621c887" />
    Configuration
</h2>

Row limits, column widths, color schemes (`dracula`, `gruvbox`, `catppuccin-mocha`, etc.) and UI visibility — all in `~/.config/squix/config.yaml`.

See [Configuration](docs/configuration.md)

---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle;" src="https://github.com/user-attachments/assets/c46a2565-a58c-472c-9393-96724d9716da" />
    Database Support
</h2>

PostgreSQL, MySQL, MariaDB, SQL Server, SQLite, Oracle, ClickHouse, Firebird

See init examples and dbeesly in [Database Support](docs/databases.md)

---

### Shell Completion

Dynamic tab completion for bash, zsh, and fish — includes your saved queries and connections.

See [Shell Completion](docs/completion.md)

---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle;" src="https://github.com/user-attachments/assets/c125a9f2-d4b6-4ec3-aef4-f52e1c8f48e8" />
    Features
</h2>

- **Query Management** — Save, organize, and execute SQL queries with parameterized support
- **TUI Table Viewer** — Vim-style navigation, in-place editing, visual selection
- **Connection Switching** — Manage multiple databases and switch instantly
- **Database Exploration** — Browse schema, visualize foreign key relationships
- **Editor Integration** — Uses `$EDITOR` for editing queries and data
- **Interactive Shell** — REPL with history, multi-line, and meta-commands

See [Features](docs/features.md) for details and examples

---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle;" src="https://github.com/user-attachments/assets/4b1425ae-7918-4a3f-b37c-41c3e443929e" />
    All Commands
</h2>

See [Commands](docs/commands.md) for the full command reference and database init examples

---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle;" src="https://github.com/user-attachments/assets/504a8488-69bf-43b4-860b-0659a6db3c69" />
    TUI Table Navigation
</h2>

See [Keybindings](docs/keybindings.md) for all navigation, editing, search, and visual mode keybindings

---

<h2>
    <img width="auto" height="24" alt="image" style="vertical-align:middle;" src="https://github.com/user-attachments/assets/432c6b41-b2e0-4326-a3cc-7b349a987bb0" />
    Roadmap
</h2>

> This project is currently in beta, please report unexpected behavior through the issues tab

### v0.3.0 - Squix 🐿️
- [x] Edit command overhaul
- [x] Delete connections with remove command
- [x] Full project rename

### v0.4.0 - Acorn 🌰
- [x] Interactive query shell (`squix shell`)
- [x] Shell autocomplete (bash, fish, zsh)
- [x] Cell search (`/`) and column header search (`f`)
- [ ] Encryption on connection username/password in config file
- [ ] Dynamic column width
- [ ] Duckdb support
- [ ] Update to bubbletea v2

---

## Contributing

We welcome contributions! Get started with detailed instructions from [CONTRIBUTING.md](CONTRIBUTING.md)

Thanks a lot to all the contributors:

<a href="https://github.com/DeprecatedLuar"><img src="https://github.com/DeprecatedLuar.png" width="40" /></a>
<a href="https://github.com/caiolandgraf"><img src="https://github.com/caiolandgraf.png" width="40" /></a>
<a href="https://github.com/g4brielklein"><img src="https://github.com/g4brielklein.png" width="40" /></a>
<a href="https://github.com/eduardofuncao"><img src="https://github.com/eduardofuncao.png" width="40" /></a>
<a href="https://github.com/udirona"><img src="https://github.com/udirona.png" width="40" /></a>
<a href="https://github.com/Leosallin"><img src="https://github.com/Leosallin.png" width="40" /></a>

## Acknowledgments

Squix wouldn't exist without the inspiration and groundwork laid by these fantastic projects:

- **[naggie/dstask](https://github.com/naggie/dstask)** - For the elegant CLI design patterns and file-based data storage approach
- **[DeprecatedLuar/better-curl-saul](https://github.com/DeprecatedLuar/better-curl-saul)** - For demonstrating a simple and genius approach to making a CLI tool
- **[dbeaver](https://github.com/dbeaver/dbeaver)** - The OG database management tool


Built with: 
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The TUI framework
- Go standard library and various database drivers

---

## License

MIT License - see [LICENSE](LICENSE) file for details

---

<div align="center">

**Made with 🐿️ by [@eduardofuncao](https://github.com/eduardofuncao)**

<img width="320" height="224" alt="Squix mascot" src="https://github.com/user-attachments/assets/f995ce07-3742-4e98-b737-bbdbf982012e" />

Previously Pam's Database Drawer, thanks to [u/marrsd](https://www.reddit.com/user/marrsd/) for suggesting the new name!

</div>
