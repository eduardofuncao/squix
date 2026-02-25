//go:build !cgo

package db

import (
	"database/sql"
	"fmt"
)

type SQLiteConnection struct {
	*BaseConnection
	db *sql.DB
}

func NewSQLiteConnection(name, connStr string) (*SQLiteConnection, error) {
	return nil, fmt.Errorf(
		"SQLite driver not available: this binary was built without CGO support. Please use a build with CGO enabled or choose a different database",
	)
}

func (oc *SQLiteConnection) Open() error {
	return fmt.Errorf("SQLite driver not available: binary built without CGO")
}

func (oc *SQLiteConnection) Ping() error {
	return fmt.Errorf("SQLite driver not available: binary built without CGO")
}

func (oc *SQLiteConnection) Close() error {
	return nil
}

func (oc *SQLiteConnection) Query(queryName string, args ...any) (any, error) {
	return nil, fmt.Errorf(
		"SQLite driver not available: binary built without CGO",
	)
}

func (oc *SQLiteConnection) ExecQuery(
	sql string,
	args ...any,
) (*sql.Rows, error) {
	return nil, fmt.Errorf(
		"SQLite driver not available: binary built without CGO",
	)
}

func (oc *SQLiteConnection) Exec(sql string, args ...any) error {
	return fmt.Errorf("SQLite driver not available: binary built without CGO")
}

func (oc *SQLiteConnection) GetTableMetadata(
	tableName string,
) (*TableMetadata, error) {
	return nil, fmt.Errorf(
		"SQLite driver not available: binary built without CGO",
	)
}

func (oc *SQLiteConnection) GetInfoSQL(infoType string) string {
	return ""
}

func (oc *SQLiteConnection) BuildUpdateStatement(
	tableName, columnName, currentValue, pkColumn, pkValue string,
) string {
	return "-- SQLite driver not available: binary built without CGO"
}

func (oc *SQLiteConnection) ApplyRowLimit(sql string, limit int) string {
	return sql
}

func (oc *SQLiteConnection) BuildDeleteStatement(
	tableName, primaryKeyCol, pkValue string,
) string {
	return "-- SQLite driver not available: binary built without CGO"
}

func (oc *SQLiteConnection) GetColumnDetails(
	tableName string,
) ([]ColumnInfo, error) {
	return nil, fmt.Errorf(
		"SQLite driver not available: binary built without CGO",
	)
}

func (oc *SQLiteConnection) BuildAddColumnSQL(
	tableName, columnName, dataType string,
	nullable bool,
	defaultValue string,
) string {
	return "-- SQLite driver not available: binary built without CGO"
}

func (oc *SQLiteConnection) BuildAlterColumnSQL(
	tableName, columnName, newDataType string,
	nullable bool,
	newDefault string,
) string {
	return "-- SQLite driver not available: binary built without CGO"
}

func (oc *SQLiteConnection) BuildRenameColumnSQL(
	tableName, oldName, newName string,
) string {
	return "-- SQLite driver not available: binary built without CGO"
}

func (oc *SQLiteConnection) BuildDropColumnSQL(
	tableName, columnName string,
) string {
	return "-- SQLite driver not available: binary built without CGO"
}

func (oc *SQLiteConnection) GetPlaceholder(paramIndex int) string {
	return "?"
}
