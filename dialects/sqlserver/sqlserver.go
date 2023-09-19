package sqlserver

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
	if isNumeric(column, columnsType) || isDatetime(column, columnsType) {
		if order == "asc" {
			return fmt.Sprintf("%s ASC", column)
		}
		return fmt.Sprintf("%s DESC", column)
	}
	if order == "asc" {
		//(CASE WHEN [Order] IS NULL THEN 0 ELSE 1 END), [Order] ASC
		return fmt.Sprintf("%s COLLATE SQL_Latin1_General_Cp1_CS_AS ASC", column)
	}
	//(CASE WHEN [Order] IS NULL THEN 1 ELSE 0 END), [Order] DESC
	return fmt.Sprintf("%s COLLATE SQL_Latin1_General_Cp1_CS_AS DESC", column)
}

func dbConfig(_ *gorm.DB) {
}
