package engine

import (
	"github.com/juaismar/go-gormssp/structs"
	"gorm.io/gorm"
)

// Simple is a main method, externally called, only return Data
func DataSimple(c Controller, conn *gorm.DB,
	table string,
	columns []structs.Data) (responseJSON structs.MessageDataTable, err error) {

	err = selectDialect(conn)
	if err != nil {
		return
	}

	responseJSON.Draw = drawNumber(c)
	myDialectFunction.DBConfig(conn)

	fieldAlias := buildType(table, conn)

	columnsType, err := initBinding(conn, "*", table, make([]structs.JoinData, 0), fieldAlias)

	// Build the SQL query string from the request
	rows, err := conn.Select("*").
		Where(filterGlobal(c, columns, columnsType, conn)).
		Where(filterIndividual(c, columns, columnsType, conn)).
		Scopes(limit(c),
			order(c, columns, columnsType)).
		Table(table).
		Rows()
	defer rows.Close()
	if err != nil {
		return
	}

	responseJSON.Data, err = dataOutput(columns, rows, columnsType)

	return
}

// Complex is a main method, externally called
func DataComplex(c Controller, conn *gorm.DB, table string, columns []structs.Data,
	whereResult []string,
	whereAll []string,
	whereJoin []structs.JoinData) (responseJSON structs.MessageDataTable, err error) {

	err = selectDialect(conn)
	if err != nil {
		return
	}

	responseJSON.Draw = drawNumber(c)
	myDialectFunction.DBConfig(conn)

	// Build the SQL query string from the request
	whereResultFlated := Flated(whereResult)
	whereAllFlated := Flated(whereAll)

	selectQuery, fieldAlias, err := buildSelectAndType(table, whereJoin, conn)
	if err != nil {
		return
	}

	columnsType, err := initBinding(conn, selectQuery, table, whereJoin, fieldAlias)
	if err != nil {
		return
	}

	rows, err := conn.Select(selectQuery).
		Where(filterGlobal(c, columns, columnsType, conn)).
		Where(filterIndividual(c, columns, columnsType, conn)).
		Scopes(
			setJoins(whereJoin),
			limit(c),
			order(c, columns, columnsType)).
		Where(whereResultFlated).
		Where(whereAllFlated).
		Table(table).
		Rows()

	if err != nil {
		return
	}
	defer rows.Close()

	responseJSON.Data, err = dataOutput(columns, rows, columnsType)
	rows.Close()

	return
}
