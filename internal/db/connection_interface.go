package db

import "database/sql"

type DatabaseConnection interface {
	Open() error
	Ping() error
	Close() error
	Query(queryName string, args ...any) (any, error)
	ExecQuery(sql string, args ...any) (*sql.Rows, error)
	Exec(sql string, args ...any) error
	GetInfoSQL(infoType string) string
	GetTables() ([]string, error)
	GetViews() ([]string, error)
	GetTableMetadata(tableName string) (*TableMetadata, error)
	GetColumnDetails(tableName string) ([]ColumnInfo, error)
	GetForeignKeys(tableName string) ([]ForeignKey, error)
	GetForeignKeysReferencingTable(tableName string) ([]ForeignKey, error)
	GetUniqueConstraints(tableName string) ([]string, error)
	BuildUpdateStatement(
		tableName, columnName, currentValue, pkColumn, pkValue string,
	) string
	BuildDeleteStatement(tableName, primaryKeyCol, pkValue string) string
	BuildAddColumnSQL(
		tableName, columnName, dataType string,
		nullable bool,
		defaultValue string,
	) string
	BuildAlterColumnSQL(
		tableName, columnName, newDataType string,
		nullable bool,
		newDefault string,
	) string
	BuildRenameColumnSQL(tableName, oldName, newName string) string
	BuildDropColumnSQL(tableName, columnName string) string
	ApplyRowLimit(sql string, limit int) string
	GetPlaceholder(paramIndex int) string

	GetName() string
	GetDbType() string
	GetConnString() string
	GetSchema() string
	GetQueries() map[string]Query
	GetLastQuery() Query

	SetSchema(string)
	SetLastQuery(Query)
	SetQueries(map[string]Query)
}
