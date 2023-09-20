package sqlserver

import (
	"database/sql"
	"fmt"
	"strings"

	dialects "github.com/juaismar/go-gormssp/dialects"
	"gorm.io/gorm"
)

func ExampleFunctions() *dialects.DialectFunctions {
	return &dialects.DialectFunctions{
		Order:    checkOrder,
		DBConfig: dbConfig,
	}
}

// Exported functions
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

// Auxiliary functions

func isNumeric(column string, columnsType []*sql.ColumnType) bool {
	for _, columnInfo := range columnsType {
		if strings.Replace(column, "\"", "", -1) == columnInfo.Name() {
			searching := columnInfo.DatabaseTypeName()
			return bindingTypesNumeric(searching, columnInfo)
		}
	}

	return false
}
func isDatetime(column string, columnsType []*sql.ColumnType) bool {
	for _, columnInfo := range columnsType {
		if strings.Replace(column, "\"", "", -1) == columnInfo.Name() {
			searching := columnInfo.DatabaseTypeName()
			return searching == "datetime" || searching == "TIMESTAMPTZ" || searching == "DATETIMEOFFSET" || searching == "DATETIME"
		}
	}

	return false
}

func bindingTypesNumeric(searching string, columnInfo *sql.ColumnType) bool {
	switch clearSearching(searching) {
	case "int", "REAL", "NUMERIC", "FLOAT":
		return true
	case "bool", "BOOL", "numeric", "BIT":
		return true
	default:
		return false
	}
}

func clearSearching(searching string) string {
	tipeElement := strings.ToLower(searching)

	switch {
	case strings.Contains(tipeElement, "varchar"):
		return "varchar"
	case strings.Contains(tipeElement, "int"):
		return "int"
	default:
		return searching
	}
}
