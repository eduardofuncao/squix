# TUI Table Navigation

When viewing query results in the TUI, you have full Vim-style navigation and editing capabilities.

## Basic Navigation

| Key | Action |
|-----|--------|
| `h`, `←` | Move left |
| `j`, `↓` | Move down |
| `k`, `↑` | Move up |
| `l`, `→` | Move right |
| `g` | Jump to first row |
| `G` | Jump to last row |
| `0`, `_`, `Home` | Jump to first column |
| `$`, `End` | Jump to last column |
| `Ctrl+u`, `PgUp` | Page up |
| `Ctrl+d`, `PgDown` | Page down |

## Data Operations

| Key | Action |
|-----|--------|
| `v` | Enter visual selection mode |
| `y` | Copy selected cell(s) to clipboard |
| `x` | Export selected cell(s) as csv, tsv, json, sql insert statement, markdown or html to clipboard |
| `Enter` | Show cell value in detail view (with JSON formatting) |
| `u` | Update current cell |
| `D` | Delete current row |
| `e` | Edit and re-run query |
| `s` | Save current query |
| `?` | Toggle keybindings help in footer |
| `q`, `Ctrl+c`, `Esc` | Quit table view |

## Search

| Key | Action |
|-----|--------|
| `/` | Search cell content |
| `n` | Jump to next cell match |
| `N` | Jump to previous cell match |
| `f` | Search column names |
| `;` | Jump to next column match |
| `,` | Jump to previous column match |

## Detail View Mode

Press `Enter` on any cell to open a detailed view that shows the full cell content. If the content is valid JSON, it will be automatically formatted with proper indentation.

**In Detail View:**

| Key | Action |
|-----|--------|
| `↑`, `↓`, `j`, `k` | Scroll through content |
| `e` | Edit cell content (opens editor with formatted JSON) |
| `q`, `Esc`, `Enter` | Close detail view |
