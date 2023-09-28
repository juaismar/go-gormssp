package structs

import (
	"database/sql"
)

type ColumnType struct {
	ColumnName    string // the field alias
	OriginalName  string // internal working name
	Type          string
	SQLColumnType *sql.ColumnType
}
