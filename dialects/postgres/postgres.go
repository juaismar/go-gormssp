package postgres

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

/*
func checkOrderDialect(column, order string, columnsType []*sql.ColumnType) string {
	const asc = "ASC NULLS FIRST"
	const desc = "DESC NULLS LAST"

	switch {
	case isSQLite(dialect) && !(isNumeric(column, columnsType) || isDatetime(column, columnsType)):
		if order == "asc" {
			return fmt.Sprintf("%s %s", column, desc)
		}
		return fmt.Sprintf("%s %s", column, asc)
	case dialect == "sqlserver":
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
	default:
		if order == "asc" {
			return fmt.Sprintf("%s %s", column, asc)
		}
		return fmt.Sprintf("%s %s", column, desc)
	}
}*/
