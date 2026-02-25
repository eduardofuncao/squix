//go:build !cgo

package db

import (
	"database/sql"
	"fmt"
)

type OracleConnection struct {
	*BaseConnection
	db *sql.DB
}

func NewOracleConnection(name, connStr string) (*OracleConnection, error) {
	return nil, fmt.Errorf(
		"Oracle driver not available: this binary was built without CGO support. Please use a build with CGO enabled or choose a different database",
	)
}

func (oc *OracleConnection) Open() error {
	return fmt.Errorf("Oracle driver not available: binary built without CGO")
}

func (oc *OracleConnection) Ping() error {
	return fmt.Errorf("Oracle driver not available: binary built without CGO")
}

func (oc *OracleConnection) Close() error {
	return nil
}

func (oc *OracleConnection) Query(queryName string, args ...any) (any, error) {
	return nil, fmt.Errorf(
		"Oracle driver not available: binary built without CGO",
	)
}

func (oc *OracleConnection) ExecQuery(
	sql string,
	args ...any,
) (*sql.Rows, error) {
	return nil, fmt.Errorf(
		"Oracle driver not available: binary built without CGO",
	)
}

func (oc *OracleConnection) Exec(sql string, args ...any) error {
	return fmt.Errorf("Oracle driver not available: binary built without CGO")
}

func (oc *OracleConnection) GetTableMetadata(
	tableName string,
) (*TableMetadata, error) {
	return nil, fmt.Errorf(
		"Oracle driver not available: binary built without CGO",
	)
}

func (oc *OracleConnection) GetInfoSQL(infoType string) string {
	return ""
}

func (oc *OracleConnection) BuildUpdateStatement(
	tableName, columnName, currentValue, pkColumn, pkValue string,
) string {
	return "-- Oracle driver not available: binary built without CGO"
}

func (oc *OracleConnection) ApplyRowLimit(sql string, limit int) string {
	return sql
}

func (oc *OracleConnection) BuildDeleteStatement(
	tableName, primaryKeyCol, pkValue string,
) string {
	return "-- Oracle driver not available: binary built without CGO"
}

func (oc *OracleConnection) GetColumnDetails(
	tableName string,
) ([]ColumnInfo, error) {
	return nil, fmt.Errorf(
		"Oracle driver not available: binary built without CGO",
	)
}

func (oc *OracleConnection) BuildAddColumnSQL(
	tableName, columnName, dataType string,
	nullable bool,
	defaultValue string,
) string {
	return "-- Oracle driver not available: binary built without CGO"
}

func (oc *OracleConnection) BuildAlterColumnSQL(
	tableName, columnName, newDataType string,
	nullable bool,
	newDefault string,
) string {
	return "-- Oracle driver not available: binary built without CGO"
}

func (oc *OracleConnection) BuildRenameColumnSQL(
	tableName, oldName, newName string,
) string {
	return "-- Oracle driver not available: binary built without CGO"
}

func (oc *OracleConnection) BuildDropColumnSQL(
	tableName, columnName string,
) string {
	return "-- Oracle driver not available: binary built without CGO"
}

func (oc *OracleConnection) GetPlaceholder(paramIndex int) string {
	return fmt.Sprintf(":%d", paramIndex)
}
