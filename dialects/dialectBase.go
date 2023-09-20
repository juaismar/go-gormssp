package dialectBase

import (
	"database/sql"
	"reflect"

	"gorm.io/gorm"
)

// Data is a line in map that link the database field with datatable field
type Data struct {
	Db        string                                                                  //name of column
	Dt        interface{}                                                             //id of column in client (int or string)
	Cs        bool                                                                    //case sensitive - optional default false
	Sf        string                                                                  //Search Function - for custom functions declared in your ddbb
	Formatter func(data interface{}, row map[string]interface{}) (interface{}, error) // - optional
}

type DialectFunctions struct {

	/* Must order properly
	1 param is column name in ddbb
	2 param must be "asc" or "desc"
	3 param is the column type array for check type
	return a sql string*/
	Order func(string, string, []*sql.ColumnType) string

	/*For configure ddbb*/
	DBConfig func(*gorm.DB)

	BindingTypesQuery func(string, string, string, *sql.ColumnType, bool, Data) (string, interface{})

	ParseData func(string, string, interface{}, reflect.Type) (interface{}, error)
}

/*
//Original order
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

/*
//Original DBConfig
func dbConfig(conn *gorm.DB) {
	if isSQLite(dialect) {
		conn.Exec("PRAGMA case_sensitive_like = ON;")
	}
}
*/

/*

func bindingTypesQuery(searching, columndb, value string, columnInfo *sql.ColumnType, isRegEx bool, column dialects.Data) (string, interface{}) {
	var fieldName = columndb
	if column.Sf != "" { //if implement custom search function
		fieldName = column.Sf
	}

	switch clearSearching(searching) {
	case "string", "TEXT", "varchar", "text":
		if isRegEx {
			return regExp(fieldName, value)
		}

		if column.Cs {
			if dialect == "sqlserver" {
				return fmt.Sprintf("%s COLLATE SQL_Latin1_General_Cp1_CS_AS LIKE ?", fieldName), "%" + value + "%"
			}
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
		fmt.Printf("(004) GORMSSP New type %v\n", columnInfo.DatabaseTypeName())
		return "", ""
	}
}

func regExp(columndb, value string) (string, string) {
	switch dialect {
	case "sqlite", "sqlite3":
		//TODO make regexp
		return fmt.Sprintf("Lower(%s) LIKE ?", columndb), "%" + strings.ToLower(value) + "%"
	case "postgres":
		return fmt.Sprintf("%s ~* ?", columndb), value
	case "sqlserver":
		//TODO make regexp
		return fmt.Sprintf("%s LIKE ?", columndb), value
	default:
		return fmt.Sprintf("%s ~* ?", columndb), value
	}
}
*/

/*

func getFieldsSearch(searching, key string, val interface{}, vType reflect.Type) (interface{}, error) {
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
*/
