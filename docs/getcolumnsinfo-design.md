# GetColumnsInfo: rich per-column metadata

Design for extending the schema explorer's columns view (`c` keybind) with
nullable / default / auto-increment information — without slowing down query
execution.

## Motivation

The explorer's `c` view currently shows:

```
column | type | pk | fk | unique
```

built from `GetTableMetadata` (columns, types, primary keys, foreign keys,
unique constraints). The most-requested missing fields are **nullability**,
**default value**, and **auto-increment/identity**.

## Why not enrich `TableMetadata`

`GetTableMetadata` runs on the **hot path** of every `squix run`:

1. `run.ExecuteSelect` → `extractMetadata` → `db.InferTableMetadata` →
   `conn.GetTableMetadata` (internal/run/executor.go)
2. `table.New` calls it twice more during TUI render (column types + FK icons)
   (internal/table/model.go)

The run path only needs `Columns`, `ColumnTypes`, `PrimaryKeys`,
`ForeignKeys`. Widening the underlying catalog SELECTs to fetch nullable /
default / comments would make every query slower — twice — for data it never
reads. Keep `GetTableMetadata` lean.

## Proposal

A dedicated interface method with a narrow contract, called **only** from the
explorer's columns action:

```go
// internal/db
type ColumnInfo struct {
    Name          string
    Type          string
    Nullable      bool
    Default       string // empty when none
    AutoIncrement bool
}

GetColumnsInfo(tableName string) ([]ColumnInfo, error)
```

- One targeted catalog query per driver (the same SQL already inside each
  driver's `GetTableMetadata` column loop, widened).
- PK/FK/unique flags stay in the view layer, cross-referenced by column name
  from the existing `GetTableMetadata` slices.
- Base connection returns an error ("not implemented") like other optional
  methods; the explorer falls back to the current view on error.

## Per-driver field mapping

| Driver     | Catalog source                | Nullable          | Default                  | AutoIncrement                     |
|------------|-------------------------------|-------------------|--------------------------|-----------------------------------|
| Postgres   | `information_schema.columns`  | `is_nullable`     | `column_default`         | `is_identity` / `column_default LIKE 'nextval%'` |
| MySQL      | `INFORMATION_SCHEMA.COLUMNS`  | `IS_NULLABLE`     | `COLUMN_DEFAULT`         | `EXTRA LIKE '%auto_increment%'`   |
| SQLite     | `PRAGMA table_info`           | `notnull = 0`     | `dflt_value` (already scanned today and discarded!) | `sql GLOB '*AUTOINCREMENT*'` on `sqlite_master.sql` |
| SQL Server | `sys.columns`                 | `is_nullable = 1` | default object resolution | `is_identity = 1`               |
| Oracle     | `all_tab_columns`             | `NULLABLE = 'Y'`  | `DATA_DEFAULT`           | `IDENTITY_COLUMN = 'YES'`         |
| DuckDB     | `information_schema.columns`  | `is_nullable`     | `column_default`         | sequence check                    |
| ClickHouse | `system.columns`              | — (always nullable unless NOT NULL type) | `default_expression` | `default_kind = 'auto_increment'` |
| Snowflake  | `INFORMATION_SCHEMA.COLUMNS` / `SHOW COLUMNS` | — | — | `AUTOINCREMENT` in DDL only |
| Firebird   | `RDB$RELATION_FIELDS` + `RDB$FIELDS` | `RDB$NULL_FLAG` (already queried via `GetInfoSQL("columns")` but discarded) | — | identity fields (`RDB$IDENTITY_TYPE`) |

Notes:
- SQLite already reads `notnull` and `dflt_value` into locals in
  `GetTableMetadata` and throws them away — that driver is a near-free win.
- Firebird already queries nullable/char-length through its internal
  `GetInfoSQL("columns")` case; same discard situation.
- Snowflake has no reliable catalog-level default/identity info; return what's
  available and leave the rest empty.

## View changes

Explorer columns view becomes:

```
column | type | nullable | default | pk | fk | unique
```

- `nullable`: `✓` or empty (or `NO` marker — pick one convention).
- `default`: raw expression text, truncated by the existing cell-width logic.
- Auto-increment rendered as part of the `type` or a dedicated marker (`↻`),
  decided at implementation time.

## Testing

- Unit-testable per driver where a live DB isn't needed (SQLite via temp file).
- Integration tests exist for sqlite already (`test/integration`); extend the
  seeded schema with NULL/DEFAULT/AUTOINCREMENT columns and assert the `c`
  path via `GetColumnsInfo` directly (TUI itself needs manual verification).

## Deferred

- Column comments: four different mechanisms (MySQL `COLUMN_COMMENT`,
  PG `col_description()`, SQL Server `sys.extended_properties`,
  Oracle `all_col_comments`). Worth doing after the core lands.
- Generated/computed columns (`is_generated`, `generation_expression`).
- Indexes artifact (name, columns, unique) as its own explorer keybind.

Origin: supersedes `info-artifacts-design.md` (removed); the per-table detail
concept lives here now since the explorer is the per-table surface.
