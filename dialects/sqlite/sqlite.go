package sqlite

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
	const asc = "ASC NULLS FIRST"
	const desc = "DESC NULLS LAST"

	if !(isNumeric(column, columnsType) || isDatetime(column, columnsType)) {
		if order == "asc" {
			return fmt.Sprintf("%s %s", column, desc)
		}
		return fmt.Sprintf("%s %s", column, asc)
	}

	if order == "asc" {
		return fmt.Sprintf("%s %s", column, asc)
	}
	return fmt.Sprintf("%s %s", column, desc)
}

func dbConfig(conn *gorm.DB) {
	conn.Exec("PRAGMA case_sensitive_like = ON;")
}
