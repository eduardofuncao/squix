package db

import (
	"fmt"
)

func CreateConnection(name, dbType, connString string) (DatabaseConnection, error) {
	switch dbType {
	case "postgres", "postgresql":
		return NewPostgresConnection(name, EncodeUserinfo(connString))
	case "mysql", "mariadb":
		return NewMySQLConnection(name, connString)
	case "sqlite", "sqlite3":
		return NewSQLiteConnection(name, connString)
	case "sqlserver", "mssql":
		return NewSQLServerConnection(name, EncodeUserinfo(connString))
	case "duckdb":
		return NewDuckDBConnection(name, connString)
	case "clickhouse":
		return NewClickHouseConnection(name, EncodeUserinfo(connString))
	case "godror", "oracle":
		return NewOracleConnection(name, EncodeUserinfo(connString))
	case "firebird", "interbase":
		return NewFirebirdConnection(name, EncodeUserinfo(connString))
	case "snowflake":
		return NewSnowflakeConnection(name, EncodeUserinfo(connString))
	default:
		return nil, fmt.Errorf("driver not implemented for %s", dbType)
	}
}
