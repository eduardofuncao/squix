package run

import "github.com/eduardofuncao/squix/internal/db"

type Flags struct {
	EditMode      bool
	LastQuery     bool
	Selector      string
	ExportFormat  string
	HideQueryName bool // suppress the query name in the TUI header
	HideQuerySQL  bool // suppress the query SQL in the TUI header
}

type ResolvedQuery struct {
	Query    db.Query
	Saveable bool // will be saved to config file
}
