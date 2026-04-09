# Commands

## Connection Management

| Command | Description | Example |
|---------|-------------|---------|
| `create <name> <type> <conn-string> [schema]` | Create new database connection | `squix create mydb postgres "postgresql://..."` |
| `switch <name>` | Switch to a different connection | `squix switch production` |
| `status` | Show current active connection | `squix status` |
| `list connections` | List all configured connections | `squix list connections` |

## Query Operations

| Command | Description | Example |
|---------|-------------|---------|
| `add <name> [sql]` | Add a new saved query | `squix add users "SELECT * FROM users"` |
| `remove <name\|id>` | Remove a saved query | `squix remove users` or `squix remove 3` |
| `list queries` | List all saved queries | `squix list queries` |
| `list queries --oneline` | lists each query in one line | `squix list -o` |
| `list queries <searchterm>` | lists queries containing search term | `squix list employees` |
| `run <name\|id\|sql>` | Execute a query | `squix run users` or `squix run 2` |
| `run` | Create and run a new query | `squix run` |
| `run --edit` | Edit query before running | `squix run users --edit` |
| `run --last`, `-l` | Re-run last executed query | `squix run --last` |
| `run --param` | run with named params | `squix run --name Squix` |
| `shell` | Interactive query REPL (alias: `repl`) | `squix shell` |


## Database Exploration

| Command | Description | Example |
|---------|-------------|---------|
| `explore` | List all tables and views in multi-column format | `squix explore` |
| `explore <table> [-l N]` | Query a table with optional row limit | `squix explore employees --limit 100` |
| `explain <table> [-d N] [-c]` | Visualize foreign key relationships | `squix explain employees --depth 2` |

## Info

| Command | Description | Example |
|---------|-------------|---------|
| `info tables` | List all tables from current schema | `squix info tables` |
| `info views` | List all views from current schema | `squix info views` |

## Configuration

| Command | Description | Example |
|---------|-------------|---------|
| `config` | Edit main configuration file | `squix config` |
| `edit` | Edit all queries for current connection | `squix edit` |
| `edit <name\|id>` | Edit a single named query | `squix edit 3` |
| `help [command]` | Show help information | `squix help run` |

---

## Database Support — Init Examples

### PostgreSQL

```bash
squix init pg-prod postgres postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable

# or connect to a specific schema:
squix init pg-prod postgres postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable schema-name
```

### MySQL / MariaDB

```bash
squix init mysql-dev mysql 'myuser:mypassword@tcp(127.0.0.1:3306)/mydb'

squix init mariadb-docker mariadb "root:MyStrongPass123@tcp(localhost:3306)/forestgrove"
```

### SQL Server


```bash
squix init sqlserver-docker sqlserver "sqlserver://sa:MyStrongPass123@localhost:1433/master"
```

### SQLite

```bash
squix init sqlite-local sqlite file:///home/eduardo/dbeesly/sqlite/mydb.sqlite
```

### Oracle

```bash
squix init oracle-stg oracle "oracle://myuser:mypassword@localhost:1521/XEPDB1"

# or connect to a specific schema:
squix init oracle-stg oracle "oracle://myuser:mypassword@localhost:1521/XEPDB1" schema-name
```

### ClickHouse

```bash
squix init clickhouse-docker clickhouse "clickhouse://myuser:mypassword@localhost:9000/forestgrove"
```

### FireBird

```bash
squix init firebird-docker firebird user:masterkey@localhost:3050//var/lib/firebird/data/the_office
```
