package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// LastQueryPath returns the path to the last-query recovery file.
func LastQueryPath() string {
	return filepath.Join(CfgPath, "last-query.sql")
}

// LoadLastQuery reads the last attempted editor-opened SQL.
// Returns ("", nil) if no file exists.
func LoadLastQuery() (string, error) {
	data, err := os.ReadFile(LastQueryPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// SaveLastQuery persists the SQL for recovery. No-op for empty input.
func SaveLastQuery(sql string) error {
	if strings.TrimSpace(sql) == "" {
		return nil
	}
	if err := os.MkdirAll(CfgPath, 0755); err != nil {
		return err
	}
	return os.WriteFile(LastQueryPath(), []byte(sql), 0600)
}

// ClearLastQuery removes the recovery file. Ignores missing-file.
func ClearLastQuery() error {
	if err := os.Remove(LastQueryPath()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
