package backup

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const pgDumpBin = "pg_dump"

var postgresFormats = map[string]FormatSpec{
	"custom": {Flag: "-Fc", Ext: "dump"},
	"plain":  {Flag: "-Fp", Ext: "sql"},
	"tar":    {Flag: "-Ft", Ext: "tar"},
}

const postgresDefaultFormat = "custom"

type postgresDumper struct{}

func newPostgresDumper() *postgresDumper {
	return &postgresDumper{}
}

func (d *postgresDumper) DefaultFormat() string {
	return postgresDefaultFormat
}

func (d *postgresDumper) Formats() map[string]FormatSpec {
	return postgresFormats
}

func (d *postgresDumper) Dump(connString, format, outPath string) error {
	if _, err := exec.LookPath(pgDumpBin); err != nil {
		return fmt.Errorf("%s not found on PATH: %w", pgDumpBin, err)
	}

	spec, ok := postgresFormats[format]
	if !ok {
		return fmt.Errorf("unknown postgres backup format %q", format)
	}

	u, err := url.Parse(connString)
	if err != nil {
		return fmt.Errorf("invalid connection string: %w", err)
	}

	password, _ := u.User.Password()

	cmd := exec.Command(pgDumpBin, spec.Flag, "-f", outPath, connString)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+password)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return nil
}
