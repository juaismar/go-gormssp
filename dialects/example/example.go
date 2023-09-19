package example

import (
	"database/sql"
	"fmt"

	SSP "github.com/juaismar/go-gormssp"
)

// OpenDB return the Database connection
func ExampleFunctions() *SSP.DialectFunctions {
	return &SSP.DialectFunctions{
		Order: checkOrder,
	}
}

func checkOrder(column, order string, columnsType []*sql.ColumnType) string {
	if order == "asc" {
		return fmt.Sprintf("%s %s", column, "ASC NULLS FIRST")
	}
	return fmt.Sprintf("%s %s", column, "DESC NULLS LAST")

}
