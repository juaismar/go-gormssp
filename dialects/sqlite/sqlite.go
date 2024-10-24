package sqlite

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/juaismar/go-gormssp/structs"
	"gorm.io/gorm"
)

func TheFunctions() *structs.DialectFunctions {
	return &structs.DialectFunctions{
		Order:             checkOrder,
		DBConfig:          dbConfig,
		BindingTypesQuery: bindingTypesQuery,
		ParseData:         parseData,
		//BindTypes:          bindTypes,
		ReservedWords:      reservedWords,
		ParseReservedField: parseReservedField,
		EscapeChar:         escapeChar,
		AliasSeparator:     aliasSeparator,
	}
}

var reservedWords = []string{}
var escapeChar = "\""
var aliasSeparator = "."

// Exported functions
func checkOrder(column, order string, columnsType []structs.ColumnType, opt map[string]interface{}) string {
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

func dbConfig(conn *gorm.DB, opt map[string]interface{}) {
	conn.Exec("PRAGMA case_sensitive_like = ON;")
}

func bindingTypesQuery(searching, columndb, value string, columnInfo structs.ColumnType, isRegEx bool, column structs.DataParsed,
	opt map[string]interface{}) (string, interface{}) {
	var fieldName = columndb

	if column.Opt != nil {
		SfField := column.Opt["SearchFunctionField"].(string)

		if SfField != "" { //if implement custom search function
			fieldName = "(" + SfField + "(" + fieldName + "))"
		}

	}

	switch clearSearching(searching) {
	case "string", "TEXT", "varchar", "text":
		if isRegEx {
			return regExp(fieldName, value)
		}

		if column.Cs {
			return fmt.Sprintf("%s LIKE ?", fieldName), "%" + value + "%"
		}
		return fmt.Sprintf("Lower(%s) LIKE ?", fieldName), "%" + strings.ToLower(value) + "%"
	case "UUID", "blob":
		if isRegEx {
			return regExp(fmt.Sprintf("CAST(%s AS TEXT)", fieldName), value)
		}
		return fmt.Sprintf("%s = ?", fieldName), value
	case "int":
		if isRegEx {
			return regExp(fmt.Sprintf("CAST(%s AS TEXT)", fieldName), value)
		}
		intval, err := strconv.Atoi(value)
		if err != nil {
			return "", ""
		}
		return fmt.Sprintf("%s = ?", fieldName), intval
	case "bool", "BOOL", "numeric", "BIT":
		if isNil(value) {
			return fieldName, nil
		}
		boolval, _ := strconv.ParseBool(value)
		return fieldName, boolval
	case "REAL", "NUMERIC", "FLOAT":
		if isRegEx {
			return regExp(fmt.Sprintf("CAST(%s AS TEXT)", fieldName), value)
		}
		fmt.Print("(005) GORMSSP WARNING: Serarching float values, float cannot be exactly equal\n")
		float64val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", ""
		}
		return fmt.Sprintf("%s = ?", fieldName), float64val
	default:
		fmt.Printf("(004) GORMSSP New type %v\n", columnInfo.Type)
		return "", ""
	}
}

func parseData(searching, key string, val interface{}, vType reflect.Type, columnInfo structs.ColumnType,
	opt map[string]interface{}) (interface{}, error) {
	switch clearSearching(searching) {
	case "string", "TEXT", "varchar", "text":
		return val.(string), nil
	case "int":
		switch vType.String() {
		case "string":
			return val.(string), nil
		default:
			return val.(int64), nil
		}
	case "NUMERIC", "REAL", "FLOAT":
		switch vType.String() {
		case "[]uint8":
			return strconv.ParseFloat(string(val.([]uint8)), 64)
		case "string":
			return strconv.ParseFloat(val.(string), 64)
		case "float64":
			return val.(float64), nil
		default:
			return val, nil
		}
	case "bool", "BOOL", "numeric", "BIT":
		switch vType.String() {
		case "int64":
			return val.(int64) == 1, nil
		case "bool":
			return val.(bool), nil
		default:
			return val, nil
		}

	case "TIMESTAMPTZ", "datetime", "DATETIMEOFFSET", "DATETIME":
		return val.(time.Time), nil
	case "UUID", "uuid", "blob":
		switch vType.String() {
		case "[]uint8":
			return string(val.([]uint8)), nil
		case "string":
			return val, nil
		default:
			return val, nil
		}
	default:
		fmt.Printf("(006) GORMSSP New type: %v for key: %v\n", searching, key)
		return val, nil
	}
}

//func bindTypes(db *gorm.DB, tableName string) (types map[string]string) {
//	return
//}

func parseReservedField(columnName string, opt map[string]interface{}) string {
	return "\"" + columnName + "\""
}

// Auxiliary functions

func isNumeric(column string, columnsType []structs.ColumnType) bool {
	for _, columnInfo := range columnsType {
		if strings.Replace(column, "\"", "", -1) == columnInfo.ColumnName {
			return bindingTypesNumeric(columnInfo.Type, columnInfo.SQLColumnType)
		}
	}

	return false
}
func isDatetime(column string, columnsType []structs.ColumnType) bool {
	for _, columnInfo := range columnsType {
		if strings.Replace(column, "\"", "", -1) == columnInfo.ColumnName {
			return columnInfo.Type == "datetime" || columnInfo.Type == "TIMESTAMPTZ" ||
				columnInfo.Type == "DATETIMEOFFSET" || columnInfo.Type == "DATETIME"
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

func regExp(columndb, value string) (string, string) {
	//TODO make regexp
	return fmt.Sprintf("Lower(%s) LIKE ?", columndb), "%" + strings.ToLower(value) + "%"
}

func isNil(val string) bool {
	valLower := strings.ToLower(val)
	return valLower == "null" || valLower == "nil" || valLower == "undefined"
}
