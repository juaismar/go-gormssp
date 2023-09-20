package postgres

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
	if order == "asc" {
		return fmt.Sprintf("%s %s", column, "ASC NULLS FIRST")
	}
	return fmt.Sprintf("%s %s", column, "DESC NULLS LAST")
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
