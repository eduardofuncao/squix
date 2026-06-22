package db

import (
	"fmt"
)

func CreateConnection(name, dbType, connString string) (DatabaseConnection, error) {
	switch dbType {
	case "postgres", "postgresql":
		return NewPostgresConnection(name, EncodeUserinfoPassword(connString))
	case "mysql", "mariadb":
		return NewMySQLConnection(name, connString)
	case "sqlite", "sqlite3":
		return NewSQLiteConnection(name, connString)
	case "sqlserver", "mssql":
		return NewSQLServerConnection(name, EncodeUserinfoPassword(connString))
	case "duckdb":
		return NewDuckDBConnection(name, connString)
	case "clickhouse":
		return NewClickHouseConnection(name, EncodeUserinfoPassword(connString))
	case "godror", "oracle":
		return NewOracleConnection(name, EncodeUserinfoPassword(connString))
	case "firebird", "interbase":
		return NewFirebirdConnection(name, EncodeUserinfoPassword(connString))
	case "snowflake":
		return NewSnowflakeConnection(name, EncodeUserinfoPassword(connString))
	default:
		return nil, fmt.Errorf("driver not implemented for %s", dbType)
	}
}
