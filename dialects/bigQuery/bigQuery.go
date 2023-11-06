package bigQuery

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/juaismar/go-gormssp/structs"
	"gorm.io/gorm"
)

func TheFunctions() *structs.DialectFunctions {
	return &structs.DialectFunctions{
		Order:              checkOrder,
		DBConfig:           dbConfig,
		BindingTypesQuery:  bindingTypesQuery,
		ParseData:          parseData,
		BindTypes:          bindTypes,
		ReservedWords:      reservedWords,
		ParseReservedField: parseReservedField,
		EscapeChar:         escapeChar,
		AliasSeparator:     aliasSeparator,
	}
}

var reservedWords = []string{}
var escapeChar = "`"
var aliasSeparator = ":"

// Exported functions
func checkOrder(column, order string, columnsType []structs.ColumnType) string {
	for _, columnInfo := range columnsType {
		if strings.Replace(column, "\"", "", -1) == columnInfo.ColumnName {
			if order == "asc" {
				return fmt.Sprintf("%s %s", columnInfo.OriginalName, "ASC NULLS FIRST")
			}
			return fmt.Sprintf("%s %s", columnInfo.OriginalName, "DESC NULLS LAST")
		}
	}
	return fmt.Sprintf("%s %s", column, "DESC NULLS LAST")
}

func dbConfig(_ *gorm.DB) {
}

func bindingTypesQuery(searching, columndb, value string, columnInfo structs.ColumnType, isRegEx bool, column structs.Data) (string, interface{}) {

	var fieldName = columndb
	if column.Sf != "" { //if implement custom search function
		fieldName = column.Sf
	}

	switch columnInfo.Type {
	case "STRING":

		//prevent SQL injection
		valueParsed := strings.Replace(value, "'", "\\'", -1)
		if column.Cs {
			return fmt.Sprintf("%s LIKE \"%%%s%%\"", fieldName, valueParsed), ""
		}
		if isRegEx {
			return regExp(columnInfo.OriginalName, valueParsed), ""
		}

		return fmt.Sprintf("LOWER(%s) LIKE \"%%%s%%\"", columnInfo.OriginalName, strings.ToLower(valueParsed)), ""

	case "INT64":
		if isRegEx {
			return regExp(fmt.Sprintf("CAST(%s AS STRING)", fieldName), value), ""
		}
		intval, err := strconv.Atoi(value)
		if err != nil {
			return "", ""
		}
		return fmt.Sprintf("%s = %d", fieldName, intval), ""
	case "FLOAT64":
		if isRegEx {
			return regExp(fmt.Sprintf("CAST(%s AS STRING)", fieldName), value), ""
		}
		float64val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", ""
		}
		return fmt.Sprintf("%s = %f", fieldName, float64val), ""
	case "BOOL":
		if isNil(value) {
			return "", ""
		}
		boolval, _ := strconv.ParseBool(value)

		searchVal := "false"
		if boolval {
			searchVal = "true"
		}

		return fmt.Sprintf("%s = %s", fieldName, searchVal), ""
	case "ARRAY<INT64>":
		return fmt.Sprintf("%s = %s", fieldName, value), ""
	default:
		fmt.Printf("(004) GORMSSP New type %v\n", columnInfo.Type)
		return "", ""
	}
}

func parseData(searching, key string, val interface{}, vType reflect.Type, columnInfo structs.ColumnType) (interface{}, error) {
	switch columnInfo.Type {
	case "STRING":
		return val.(string), nil
	case "INT64":
		return val.(int64), nil
	case "FLOAT64":
		return val.(float64), nil
	case "BOOL":
		return val.(bool), nil
	default:
		return val, nil
	}
}

func bindTypes(db *gorm.DB, tableName string) (types map[string]string) {
	types = make(map[string]string)
	rows, _ := db.Raw("SELECT column_name, data_type " +
		"FROM `" + os.Getenv("LW_DATASET") + ".INFORMATION_SCHEMA.COLUMNS` " +
		"WHERE table_name = '" + tableName + "' " +
		"ORDER BY ordinal_position").
		Rows()

	for rows.Next() {
		var column_name, data_type string
		rows.Scan(&column_name, &data_type)
		types[column_name] = data_type
	}

	return
}

func parseReservedField(columnName string) string {
	return "`" + columnName + "`"
}

// Auxiliary functions

func isNil(val string) bool {
	valLower := strings.ToLower(val)
	return valLower == "null" || valLower == "nil" || valLower == "undefined"
}

func regExp(columndb, value string) string {
	return fmt.Sprintf("REGEXP_CONTAINS(%s,\"%s\") ", columndb, value)
}