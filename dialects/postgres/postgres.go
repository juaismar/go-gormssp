package postgres

import (
	"database/sql"
	"fmt"

	dialects "github.com/juaismar/go-gormssp/dialects"
	"gorm.io/gorm"
)

// OpenDB return the Database connection
func ExampleFunctions() *dialects.DialectFunctions {
	return &dialects.DialectFunctions{
		Order:    checkOrder,
		DBConfig: dbConfig,
	}
}

func checkOrder(column, order string, columnsType []*sql.ColumnType) string {
	if order == "asc" {
		return fmt.Sprintf("%s %s", column, "ASC NULLS FIRST")
	}
	return fmt.Sprintf("%s %s", column, "DESC NULLS LAST")
}

func dbConfig(_ *gorm.DB) {
}
