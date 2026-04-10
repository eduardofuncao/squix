package run

import "github.com/eduardofuncao/squix/internal/db"

type Flags struct {
	EditMode     bool
	LastQuery    bool
	Selector     string
	ExportFormat string
}

type ResolvedQuery struct {
	Query    db.Query
	Saveable bool // will be saved to config file
}
