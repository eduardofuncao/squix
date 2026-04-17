package db

import (
	"crypto/rsa"
	"database/sql"
	"encoding/pem"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
)

type SnowflakeConnection struct {
	*BaseConnection
	db        *sql.DB
	warehouse string
	role      string
}

func NewSnowflakeConnection(name, connStr string) (*SnowflakeConnection, error) {
	bc := &BaseConnection{
		Name:       name,
		DbType:     "snowflake",
		ConnString: connStr,
	}

	conn := &SnowflakeConnection{BaseConnection: bc}

	parsedURL, err := url.Parse(connStr)
	if err == nil {
		q := parsedURL.Query()
		if schema := q.Get("schema"); schema != "" {
			bc.Schema = schema
		}
		conn.warehouse = q.Get("warehouse")
		conn.role = q.Get("role")
	}

	return conn, nil
}

func (s *SnowflakeConnection) Open() error {
	db, err := s.openDB()
	if err != nil {
		return err
	}
	s.db = db

	// Force connection so auth errors surface before session commands.
	if err = s.db.Ping(); err != nil {
		s.db.Close()
		return fmt.Errorf("failed to connect to snowflake: %w", err)
	}

	if s.warehouse != "" {
		if _, err = s.db.Exec(fmt.Sprintf("USE WAREHOUSE %s", s.warehouse)); err != nil {
			s.db.Close()
			return fmt.Errorf("failed to set warehouse to '%s': %w", s.warehouse, err)
		}
	}

	if s.role != "" {
		if _, err = s.db.Exec(fmt.Sprintf("USE ROLE %s", s.role)); err != nil {
			s.db.Close()
			return fmt.Errorf("failed to set role to '%s': %w", s.role, err)
		}
	}

	if s.Schema != "" {
		if _, err = s.db.Exec(fmt.Sprintf("USE SCHEMA %s", s.Schema)); err != nil {
			s.db.Close()
			return fmt.Errorf("failed to set schema to '%s': %w", s.Schema, err)
		}
	}

	return nil
}

// openDB parses privateKeyFile from the DSN manually; gosnowflake key-file reading is unreliable.
func (s *SnowflakeConnection) openDB() (*sql.DB, error) {
	parsedURL, err := url.Parse(s.ConnString)
	if err != nil {
		return sql.Open("snowflake", s.ConnString)
	}

	q := parsedURL.Query()
	keyFile := q.Get("privateKeyFile")
	if keyFile == "" {
		return sql.Open("snowflake", s.ConnString)
	}

	privateKey, err := loadRSAPrivateKey(keyFile, q.Get("privateKeyPassphrase"))
	if err != nil {
		return nil, fmt.Errorf("failed to load private key from '%s': %w", keyFile, err)
	}

	cfg := &gosnowflake.Config{
		Account:       parsedURL.Host,
		User:          parsedURL.User.Username(),
		Database:      strings.TrimPrefix(parsedURL.Path, "/"),
		Schema:        q.Get("schema"),
		Warehouse:     q.Get("warehouse"),
		Role:          q.Get("role"),
		Authenticator: gosnowflake.AuthTypeJwt,
		PrivateKey:    privateKey,
	}

	dsn, err := gosnowflake.DSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build snowflake DSN: %w", err)
	}

	return sql.Open("snowflake", dsn)
}

func loadRSAPrivateKey(path, passphrase string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in %s", path)
	}

	rsaKey, err := pkcs8.ParsePKCS8PrivateKeyRSA(block.Bytes, []byte(passphrase))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return rsaKey, nil
}

func (s *SnowflakeConnection) Ping() error {
	if s.db == nil {
		return fmt.Errorf("database is not open")
	}
	return s.db.Ping()
}

func (s *SnowflakeConnection) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SnowflakeConnection) Query(queryName string, args ...any) (any, error) {
	query, exists := s.Queries[queryName]
	if !exists {
		return nil, fmt.Errorf("query not found: %s", queryName)
	}
	return s.db.Query(query.SQL, args...)
}

func (s *SnowflakeConnection) ExecQuery(sql string, args ...any) (*sql.Rows, error) {
	return s.db.Query(sql, args...)
}

func (s *SnowflakeConnection) Exec(sql string, args ...any) error {
	_, err := s.db.Exec(sql, args...)
	return err
}

func (s *SnowflakeConnection) GetTableMetadata(tableName string) (*TableMetadata, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database is not open")
	}

	metadata := &TableMetadata{
		TableName: tableName,
	}

	pkQuery := `
		SELECT kcu.COLUMN_NAME
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
		JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
			AND tc.TABLE_SCHEMA = kcu.TABLE_SCHEMA
			AND tc.TABLE_NAME = kcu.TABLE_NAME
		WHERE tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
		  AND tc.TABLE_NAME = ?
		ORDER BY kcu.ORDINAL_POSITION
	`
	pkRows, err := s.db.Query(pkQuery, strings.ToUpper(tableName))
	if err == nil {
		defer pkRows.Close()
		for pkRows.Next() {
			var col string
			if err := pkRows.Scan(&col); err == nil {
				metadata.PrimaryKeys = append(metadata.PrimaryKeys, col)
			}
		}
	}

	colQuery := `
		SELECT COLUMN_NAME, DATA_TYPE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`
	colRows, err := s.db.Query(colQuery, strings.ToUpper(tableName))
	if err != nil {
		return metadata, fmt.Errorf("failed to query snowflake column metadata: %w", err)
	}
	defer colRows.Close()

	for colRows.Next() {
		var colName, colType string
		if err := colRows.Scan(&colName, &colType); err == nil {
			metadata.Columns = append(metadata.Columns, colName)
			metadata.ColumnTypes = append(metadata.ColumnTypes, colType)
		}
	}

	fks, err := s.GetForeignKeys(tableName)
	if err == nil {
		metadata.ForeignKeys = fks
	} else {
		metadata.ForeignKeys = []ForeignKey{}
	}

	return metadata, nil
}

func (s *SnowflakeConnection) GetForeignKeys(tableName string) ([]ForeignKey, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database is not open")
	}

	query := `
		SELECT
			kcu.COLUMN_NAME,
			ccu.TABLE_NAME  AS REFERENCED_TABLE,
			ccu.COLUMN_NAME AS REFERENCED_COLUMN
		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
		JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON rc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
			AND rc.CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA
		JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE ccu
			ON rc.UNIQUE_CONSTRAINT_NAME = ccu.CONSTRAINT_NAME
			AND rc.UNIQUE_CONSTRAINT_SCHEMA = ccu.CONSTRAINT_SCHEMA
		WHERE kcu.TABLE_NAME = ?
		ORDER BY kcu.COLUMN_NAME
	`

	rows, err := s.db.Query(query, strings.ToUpper(tableName))
	if err != nil {
		// Snowflake FKs are unenforced; metadata may be absent — return empty gracefully
		return []ForeignKey{}, nil
	}
	defer rows.Close()

	var fks []ForeignKey
	for rows.Next() {
		var fk ForeignKey
		if err := rows.Scan(&fk.Column, &fk.ReferencedTable, &fk.ReferencedColumn); err == nil {
			fks = append(fks, fk)
		}
	}

	return fks, nil
}

func (s *SnowflakeConnection) GetForeignKeysReferencingTable(tableName string) ([]ForeignKey, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database is not open")
	}

	query := `
		SELECT
			kcu.COLUMN_NAME,
			kcu.TABLE_NAME  AS REFERENCING_TABLE,
			ccu.COLUMN_NAME AS REFERENCED_COLUMN
		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
		JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON rc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
			AND rc.CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA
		JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE ccu
			ON rc.UNIQUE_CONSTRAINT_NAME = ccu.CONSTRAINT_NAME
			AND rc.UNIQUE_CONSTRAINT_SCHEMA = ccu.CONSTRAINT_SCHEMA
		WHERE ccu.TABLE_NAME = ?
		ORDER BY kcu.TABLE_NAME, kcu.COLUMN_NAME
	`

	rows, err := s.db.Query(query, strings.ToUpper(tableName))
	if err != nil {
		return []ForeignKey{}, nil
	}
	defer rows.Close()

	var fks []ForeignKey
	for rows.Next() {
		var fk ForeignKey
		if err := rows.Scan(&fk.Column, &fk.ReferencedTable, &fk.ReferencedColumn); err == nil {
			fks = append(fks, fk)
		}
	}

	return fks, nil
}

func (s *SnowflakeConnection) GetUniqueConstraints(tableName string) ([]string, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database is not open")
	}

	query := `
		SELECT kcu.COLUMN_NAME
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
		JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
			AND tc.TABLE_SCHEMA = kcu.TABLE_SCHEMA
			AND tc.TABLE_NAME = kcu.TABLE_NAME
		WHERE tc.CONSTRAINT_TYPE = 'UNIQUE'
		  AND tc.TABLE_NAME = ?
		ORDER BY kcu.COLUMN_NAME
	`

	rows, err := s.db.Query(query, strings.ToUpper(tableName))
	if err != nil {
		return []string{}, nil
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err == nil {
			cols = append(cols, col)
		}
	}

	return cols, nil
}

func (s *SnowflakeConnection) GetInfoSQL(infoType string) string {
	schema := s.Schema
	var schemaFilter string
	if schema != "" {
		schemaFilter = fmt.Sprintf("AND TABLE_SCHEMA = '%s'", strings.ToUpper(schema))
	}

	switch infoType {
	case "tables":
		return fmt.Sprintf(`SELECT
			TABLE_SCHEMA  AS schema,
			TABLE_NAME    AS name,
			TABLE_OWNER   AS owner
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_TYPE = 'BASE TABLE'
		  %s
		ORDER BY TABLE_SCHEMA, TABLE_NAME`, schemaFilter)
	case "views":
		return fmt.Sprintf(`SELECT
			TABLE_SCHEMA  AS schema,
			TABLE_NAME    AS name,
			TABLE_OWNER   AS owner
		FROM INFORMATION_SCHEMA.VIEWS
		WHERE 1=1
		  %s
		ORDER BY TABLE_SCHEMA, TABLE_NAME`, schemaFilter)
	default:
		return ""
	}
}

func (s *SnowflakeConnection) GetTables() ([]string, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database is not open")
	}

	query := `
		SELECT TABLE_NAME
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_NAME
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			tables = append(tables, name)
		}
	}

	return tables, nil
}

func (s *SnowflakeConnection) GetViews() ([]string, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database is not open")
	}

	query := `
		SELECT TABLE_NAME
		FROM INFORMATION_SCHEMA.VIEWS
		ORDER BY TABLE_NAME
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query views: %w", err)
	}
	defer rows.Close()

	var views []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			views = append(views, name)
		}
	}

	return views, nil
}

func (s *SnowflakeConnection) BuildUpdateStatement(
	tableName, columnName, currentValue, pkColumn, pkValue string,
) string {
	escapedValue := strings.ReplaceAll(currentValue, "'", "''")

	if pkColumn != "" && pkValue != "" {
		escapedPkValue := strings.ReplaceAll(pkValue, "'", "''")
		return fmt.Sprintf(
			"-- Snowflake UPDATE statement\nUPDATE %s\nSET %s = '%s'\nWHERE %s = '%s';",
			tableName, columnName, escapedValue, pkColumn, escapedPkValue,
		)
	}

	return fmt.Sprintf(
		"-- Snowflake UPDATE statement\n-- No primary key specified. Edit WHERE clause manually.\nUPDATE %s\nSET %s = '%s'\nWHERE <condition>;",
		tableName, columnName, escapedValue,
	)
}

func (s *SnowflakeConnection) BuildDeleteStatement(tableName, primaryKeyCol, pkValue string) string {
	escapedPkValue := strings.ReplaceAll(pkValue, "'", "''")
	return fmt.Sprintf(
		"DELETE FROM %s\nWHERE %s = '%s';",
		tableName, primaryKeyCol, escapedPkValue,
	)
}

func (s *SnowflakeConnection) GetPlaceholder(paramIndex int) string {
	return "?"
}

func (s *SnowflakeConnection) ApplyRowLimit(sql string, limit int) string {
	trimmedSQL := strings.ToUpper(strings.TrimSpace(sql))

	if !strings.HasPrefix(trimmedSQL, "SELECT") {
		return sql
	}

	if strings.Contains(trimmedSQL, " LIMIT ") {
		return sql
	}

	trimmed := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(sql), ";"))
	return fmt.Sprintf("%s\nLIMIT %d", trimmed, limit)
}
