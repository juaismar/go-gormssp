package engine

import (
	"github.com/juaismar/go-gormssp/structs"
	"gorm.io/gorm"
)

// Simple is a main method, externally called, only return Data
func DataSimple(c Controller, conn *gorm.DB,
	table string,
	columns []structs.Data) (responseJSON structs.MessageDataTable, err error) {

	err = SelectDialect(conn)
	if err != nil {
		return
	}

	responseJSON.Draw = DrawNumber(c)
	MyDialectFunction.DBConfig(conn)

	fieldAlias := BuildType(table, conn)

	columnsType, err := InitBinding(conn, "*", table, make([]structs.JoinData, 0), fieldAlias)

	// Build the SQL query string from the request
	rows, err := conn.Select("*").
		Where(FilterGlobal(c, columns, columnsType, conn)).
		Where(FilterIndividual(c, columns, columnsType, conn)).
		Scopes(Limit(c),
			Order(c, columns, columnsType)).
		Table(table).
		Rows()
	defer rows.Close()
	if err != nil {
		return
	}

	responseJSON.Data, err = DataOutput(columns, rows, columnsType)

	return
}

// Complex is a main method, externally called
func DataComplex(c Controller, conn *gorm.DB, table string, columns []structs.Data,
	whereResult []string,
	whereAll []string,
	whereJoin []structs.JoinData) (responseJSON structs.MessageDataTable, err error) {

	err = SelectDialect(conn)
	if err != nil {
		return
	}

	responseJSON.Draw = DrawNumber(c)
	MyDialectFunction.DBConfig(conn)

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

	rows, err := conn.Select(selectQuery).
		Where(FilterGlobal(c, columns, columnsType, conn)).
		Where(FilterIndividual(c, columns, columnsType, conn)).
		Scopes(
			SetJoins(whereJoin),
			Limit(c),
			Order(c, columns, columnsType)).
		Where(whereResultFlated).
		Where(whereAllFlated).
		Table(table).
		Rows()

	if err != nil {
		return
	}
	defer rows.Close()

	responseJSON.Data, err = DataOutput(columns, rows, columnsType)
	rows.Close()

	return
}
