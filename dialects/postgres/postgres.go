package postgres

import (
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
func checkOrder(column, order string, columnsType []structs.ColumnType,
	opt map[string]interface{}) string {
	if order == "asc" {
		return fmt.Sprintf("%s %s", column, "ASC NULLS FIRST")
	}
	return fmt.Sprintf("%s %s", column, "DESC NULLS LAST")
}

func dbConfig(_ *gorm.DB, opt map[string]interface{}) {
}

func bindingTypesQuery(searching, columndb, value string, columnInfo structs.ColumnType, isRegEx bool, column structs.DataParsed,
	opt map[string]interface{}) (string, interface{}) {
	var fieldName = columndb
	SfValue := ""
	if column.Opt != nil {
		SfValue = column.Opt["SearchFunctionValue"].(string)
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
			parsedValue := "%" + value + "%"
			implementSearchFunctionValue(SfValue, &parsedValue)

			return fmt.Sprintf("%s LIKE ?", fieldName), parsedValue
		}
		parsedValue := "%" + strings.ToLower(value) + "%"
		implementSearchFunctionValue(SfValue, &parsedValue)
		return fmt.Sprintf("Lower(%s) LIKE ?", fieldName), parsedValue
	case "UUID", "blob":
		parsedValue := value
		implementSearchFunctionValue(SfValue, &parsedValue)
		if isRegEx {
			return regExp(fmt.Sprintf("CAST(%s AS TEXT)", fieldName), parsedValue)
		}
		return fmt.Sprintf("%s = ?", fieldName), parsedValue
	case "int":
		parsedValue := value
		implementSearchFunctionValue(SfValue, &parsedValue)
		if isRegEx {

			return regExp(fmt.Sprintf("CAST(%s AS TEXT)", fieldName), parsedValue)
		}
		intval, err := strconv.Atoi(parsedValue)
		if err != nil {
			return "", ""
		}

		return fmt.Sprintf("%s = ?", fieldName), intval
	case "bool", "BOOL", "numeric", "BIT":

		parsedValue := value
		implementSearchFunctionValue(SfValue, &parsedValue)
		if isNil(parsedValue) {
			return fieldName, nil
		}
		boolval, _ := strconv.ParseBool(value)
		return fieldName, boolval
	case "REAL", "NUMERIC", "FLOAT":

		parsedValue := value
		implementSearchFunctionValue(SfValue, &parsedValue)
		if isRegEx {
			return regExp(fmt.Sprintf("CAST(%s AS TEXT)", fieldName), parsedValue)
		}
		fmt.Print("(005) GORMSSP WARNING: Serarching float values, float cannot be exactly equal\n")
		float64val, err := strconv.ParseFloat(parsedValue, 64)
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
	return fmt.Sprintf("%s ~* ?", columndb), value
}

func isNil(val string) bool {
	valLower := strings.ToLower(val)
	return valLower == "null" || valLower == "nil" || valLower == "undefined"
}

func implementSearchFunctionValue(searchFunction string, value *string) {
	if searchFunction == "" {
		return
	}
	*value = "(" + searchFunction + "(" + *value + "%))"
}
