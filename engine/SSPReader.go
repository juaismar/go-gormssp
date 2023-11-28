package engine

import (
	"database/sql"

	"github.com/juaismar/go-gormssp/structs"
	"gorm.io/gorm"
)

// Simple is a main method, externally called, only return Data
func ReaderSimple(c Controller, conn *gorm.DB,
	table string,
	columns []structs.Data,
	opt map[string]interface{}) (rows *sql.Rows, err error) {

	Opt = opt
	parsedColumns, err := PreprocessDataColums(columns)
	if err != nil {
		return
	}

	err = SelectDialect(conn)
	if err != nil {
		return
	}

	MyDialectFunction.DBConfig(conn, Opt)

	fieldAlias := BuildType(table, conn)

	columnsType, err := InitBinding(conn, "*", table, make([]structs.JoinData, 0), fieldAlias)

	// Build the SQL query string from the request
	rows, err = conn.Select("*").
		Where(FilterGlobal(c, parsedColumns, columnsType, conn)).
		Where(FilterIndividual(c, parsedColumns, columnsType, conn)).
		Scopes(Limit(c),
			Order(c, parsedColumns, columnsType)).
		Table(table).
		Rows()

	return
}

// Complex is a main method, externally called
func ReaderComplex(c Controller, conn *gorm.DB, table string, columns []structs.Data,
	whereResult []string,
	whereAll []string,
	whereJoin []structs.JoinData,
	opt map[string]interface{}) (rows *sql.Rows, err error) {

	Opt = opt
	parsedColumns, err := PreprocessDataColums(columns)
	if err != nil {
		return
	}

	err = SelectDialect(conn)
	if err != nil {
		return
	}

	MyDialectFunction.DBConfig(conn, Opt)

	// Build the SQL query string from the request
	whereResultFlated := Flated(whereResult)
	whereAllFlated := Flated(whereAll)

	selectQuery, fieldAlias, err := BuildSelectAndType(table, whereJoin, conn)
	if err != nil {
		return
	}

	columnsType, err := InitBinding(conn, selectQuery, table, whereJoin, fieldAlias)
	if err != nil {
		return
	}

	rows, err = conn.Select(selectQuery).
		Where(FilterGlobal(c, parsedColumns, columnsType, conn)).
		Where(FilterIndividual(c, parsedColumns, columnsType, conn)).
		Scopes(
			SetJoins(whereJoin),
			Limit(c),
			Order(c, parsedColumns, columnsType)).
		Where(whereResultFlated).
		Where(whereAllFlated).
		Table(table).
		Rows()

	return
}
