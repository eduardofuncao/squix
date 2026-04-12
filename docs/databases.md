# Database Support
To create and save a database connection to be used with other `squix` commands, use the `init` command
```bash
init <name> <type> <conn-string> [schema]
```
If the <type> is ommited, squix will try to infer the database type from the connection string

## Init Examples

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

### DuckDB

> Requires CGO (included by default; excluded by building with `CGO_ENABLED=0`).

```bash
squix init duckdb-local duckdb /path/to/mydb.db
```

DuckDB can also query CSV and JSON files directly:

```bash
squix init duckdb-local duckdb employees.json
squix init duckdb-local duckdb employees.csv
```

For json and csv, the data will be available to query as a view, with the same name as the file
(eg. `squix run "select * from employees" from the employees.json file`)

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

---

## 🐝 Dbeesly

To run containerized test database servers for all supported databases, use the sister project [dbeesly](https://github.com/eduardofuncao/dbeesly)

<img width="879" height="571" alt="image" src="https://github.com/user-attachments/assets/c0a131eb-ea95-4523-86ac-cd00a561a5e0" />
