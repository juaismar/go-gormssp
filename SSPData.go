package ssp

import (
	"gorm.io/gorm"
)

// Simple is a main method, externally called, only return Data
func DataSimple(c Controller, conn *gorm.DB,
	table string,
	columns []Data) (responseJSON MessageDataTable, err error) {

	dialect = conn.Dialector.Name()

	responseJSON.Draw = drawNumber(c)
	dbConfig(conn)

	columnsType, err := initBinding(conn, "*", table, make(map[string]string, 0))

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

	responseJSON.Data, err = dataOutput(columns, rows)

	return
}

// Complex is a main method, externally called
func DataComplex(c Controller, conn *gorm.DB, table string, columns []Data,
	whereResult []string,
	whereAll []string,
	whereJoin map[string]string) (responseJSON MessageDataTable, err error) {

	dialect = conn.Dialector.Name()

	responseJSON.Draw = drawNumber(c)
	dbConfig(conn)

	// Build the SQL query string from the request
	whereResultFlated := flated(whereResult)
	whereAllFlated := flated(whereAll)

	selectQuery, err := buildSelect(table, whereJoin, conn)
	if err != nil {
		return
	}

	columnsType, err := initBinding(conn, selectQuery, table, whereJoin)
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

	responseJSON.Data, err = dataOutput(columns, rows)
	rows.Close()

	return
}
