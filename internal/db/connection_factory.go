package db

import "fmt"

type driverCtor func(name, connString string) (DatabaseConnection, error)

// Excluded drivers never appear here: each <driver>_connection.go registers
// itself from init() behind a build tag.
var builtDrivers = map[string]driverCtor{}

func registerDriver(ctor driverCtor, aliases ...string) {
	for _, alias := range aliases {
		builtDrivers[alias] = ctor
	}
}

var canonicalTypes = []string{
	"postgres",
	"mysql",
	"sqlite",
	"sqlserver",
	"clickhouse",
	"oracle",
	"firebird",
	"duckdb",
	"snowflake",
}

var excludedByTag = map[string]string{
	"postgres":   "nopostgres",
	"postgresql": "nopostgres",
	"mysql":      "nomysql",
	"mariadb":    "nomysql",
	"sqlite":     "nosqlite",
	"sqlite3":    "nosqlite",
	"sqlserver":  "nosqlserver",
	"mssql":      "nosqlserver",
	"clickhouse": "noclickhouse",
	"oracle":     "noracle",
	"godror":     "noracle",
	"firebird":   "nofirebird",
	"interbase":  "nofirebird",
	"duckdb":     "noduckdb",
	"snowflake":  "nosnowflake",
}

func CreateConnection(name, dbType, connString string) (DatabaseConnection, error) {
	if ctor, ok := builtDrivers[dbType]; ok {
		return ctor(name, connString)
	}
	if _, known := excludedByTag[dbType]; known {
		return nil, NotIncludedError(dbType)
	}
	return nil, fmt.Errorf("driver not implemented for %s", dbType)
}

func NotIncludedError(dbType string) error {
	return fmt.Errorf(
		"%s driver is not included in this build (excluded via -tags %s); use a full squix build or rebuild without the tag",
		dbType, excludedByTag[dbType],
	)
}

func KnownDBTypes() []string {
	return canonicalTypes
}

func IsDriverBuilt(dbType string) bool {
	_, ok := builtDrivers[dbType]
	return ok
}

func DriverBuildTag(dbType string) string {
	return excludedByTag[dbType]
}
