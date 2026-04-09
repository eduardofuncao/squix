# Configuration

Squix stores its configuration at `~/.config/squix/config.yaml`.

## Row Limit `default_row_limit: 1000`
All queries are automatically limited to prevent fetching massive result sets. Configure via `default_row_limit` in config or use explicit `LIMIT` in your SQL queries.

## Column Width `default_column_width: 15`
The width for all columns in the table TUI is fixed to a constant size, which can be configured through `default_column_width` in the config file. There are plans to make the column widths flexible in future versions.

## Color Schemes `color_scheme: "default"`
Customize the terminal UI colors with built-in schemes:

**Available schemes:**
`default`, `dracula`, `gruvbox`, `solarized`, `nord`, `monokai`
`black-metal`, `black-metal-gorgoroth`, `vesper`, `catppuccin-mocha`, `tokyo-night`, `rose-pine`, `terracotta`

Each scheme uses a 7-color palette: Primary (titles, headers), Success (success messages), Error (errors), Normal (table data), Muted (borders, help text), Highlight (selected backgrounds), Accent (keywords, strings).

## UI Visibility `ui_visibility`

Control which UI components are displayed in the table view:

```yaml
ui_visibility:
  query_name: true          # Show query name header
  query_sql: true           # Show SQL query display
  type_display: true        # Show column type indicators
  key_icons: true           # Show primary key (⚿) and foreign key (⚭) icons
  footer_cell_content: true # Show current cell preview in footer
  footer_stats: true        # Show row/col count and position in footer
  footer_keymaps: true      # Show keybindings help in footer
```

**Tip:** Press `?` in the table view to toggle the keymaps help on/off.
