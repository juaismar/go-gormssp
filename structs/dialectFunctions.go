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
	return a sql string*/
	Order func(string, string, []ColumnType) string

	/* Optional - can be empty function
	For configure ddbb*/
	DBConfig func(*gorm.DB)

	/* Mandatory
	Build a query for a field
	1 the type of field
	2 the name of field
	3 value to search
	4 aditional ddbb raw column info
	5 if must exect a regular expression
	6 Data info if needed

	return 1 a query "?" is for search value
	retrun 2 the search value
	*/
	BindingTypesQuery func(string, string, string, ColumnType, bool, Data) (string, interface{})

	/* Mandatory
	Parse a field responded to be return Numbers to int...
	1 type of field in ddbb
	2 field name in datatable
	3 value to be parsed
	4 type of field raw

	return parsed value
	retrun a error to return
	*/
	ParseData func(string, string, interface{}, reflect.Type, ColumnType) (interface{}, error)

	/* Optional - can be empty function
	Get field types
	1 ddbb
	2 tablename
	return [field name , type]*/
	BindTypes func(*gorm.DB, string) map[string]string

	/* Mandatory
	Parse a field reserved.
	1 fieldname
	return parsed fieldname
	*/
	ParseReservedField func(string) string

	/*Words to be scaped*/
	ReservedWords []string

	/*Character to scape names ej '\"'*/
	EscapeChar string

	/*Character for separate tableName from FieldName ej "."*/
	AliasSeparator string
}
