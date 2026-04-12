//go:build !cgo

package db

import "fmt"

type DuckDBConnection struct {
	*BaseConnection
}

func NewDuckDBConnection(name, connStr string) (*DuckDBConnection, error) {
	return nil, fmt.Errorf("duckdb driver not available: build with CGO_ENABLED=1 to enable")
}

func (d *DuckDBConnection) GetUniqueConstraints(tableName string) ([]string, error) {
	return nil, fmt.Errorf("duckdb driver not available")
}
