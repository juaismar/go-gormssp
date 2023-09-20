package dialectBase

import (
	"database/sql"

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

	bindingTypesQuery func(string, string, string, *sql.ColumnType, bool, Data) (string, interface{})
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
