package structs

import (
	"reflect"

	"gorm.io/gorm"
)

type DialectFunctions struct {

	/* Mandatory
	Must order properly
	1 param is column name in ddbb
	2 param must be "asc" or "desc"
	3 param is the column type array for check type
	4 optional params, can be nil

	return a sql string*/
	Order func(string, string, []ColumnType, map[string]interface{}) string

	/* Optional - can be empty function
	For configure ddbb
	1 gorm connection
	2 optional params, can be nil*/
	DBConfig func(*gorm.DB, map[string]interface{})

	/* Mandatory
	Build a query for a field
	1 the type of field
	2 the name of field
	3 value to search
	4 aditional ddbb raw column info
	5 if must exect a regular expression
	6 Data info if needed
	7 optional params, can be nil

	return 1 a query "?" is for search value
	retrun 2 the search value
	*/
	BindingTypesQuery func(string, string, string, ColumnType, bool, DataParsed, map[string]interface{}) (string, interface{})

	/* Mandatory
	Parse a field responded to be return Numbers to int...
	1 type of field in ddbb
	2 field name in datatable
	3 value to be parsed
	4 type of field raw
	5 optional params, can be nil

	return parsed value
	retrun a error to return
	*/
	ParseData func(string, string, interface{}, reflect.Type, ColumnType, map[string]interface{}) (interface{}, error)

	/* Optional - can be empty function
	Get field types
	1 ddbb
	2 tablename
	3 optional params, can be nil

	return [field name , type]*/
	BindTypes func(*gorm.DB, string, map[string]interface{}) map[string]string

	/* Mandatory
	Parse a field reserved.
	1 fieldname
	2 optional params, can be nil

	return parsed fieldname
	*/
	ParseReservedField func(string, map[string]interface{}) string

	/* Optional - Words to be scaped*/
	ReservedWords []string

	/* Mandatory - Character to scape names ej '\"'*/
	EscapeChar string

	/* Mandatory - Character for separate tableName from FieldName ej "."*/
	AliasSeparator string
}
