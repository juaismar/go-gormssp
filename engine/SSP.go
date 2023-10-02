package engine

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/juaismar/go-gormssp/dialects/bigQuery"
	"github.com/juaismar/go-gormssp/dialects/postgres"
	"github.com/juaismar/go-gormssp/dialects/sqlite"
	"github.com/juaismar/go-gormssp/dialects/sqlserver"
	"github.com/juaismar/go-gormssp/structs"

	"gorm.io/gorm"
)

var myDialectFunction *structs.DialectFunctions

// Controller emulate the beego controller
type Controller interface {
	GetString(string, ...string) string
}

// Simple is a main method, externally called
func Simple(c Controller, conn *gorm.DB,
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
	if err != nil {
		return
	}

	//search in DDBB recordsFiltered
	err = conn.
		Where(filterGlobal(c, columns, columnsType, conn)).
		Where(filterIndividual(c, columns, columnsType, conn)).
		Table(table).Count(&responseJSON.RecordsFiltered).Error
	if err != nil {
		return
	}

	//search in DDBB recordsTotal
	err = conn.Table(table).Count(&responseJSON.RecordsTotal).Error

	return
}

// Complex is a main method, externally called
func Complex(c Controller, conn *gorm.DB, table string, columns []structs.Data,
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
	if err != nil {
		return
	}

	err = conn.
		Where(filterGlobal(c, columns, columnsType, conn)).
		Where(filterIndividual(c, columns, columnsType, conn)).
		Scopes(
			setJoins(whereJoin),
		).
		Where(whereResultFlated).
		Where(whereAllFlated).
		Table(table).
		Count(&responseJSON.RecordsFiltered).Error
	if err != nil {
		return
	}

	err = conn.Table(table).
		Scopes(setJoins(whereJoin)).
		Where(whereAllFlated).Count(&responseJSON.RecordsTotal).Error

	return
}

func selectDialect(conn *gorm.DB) (err error) {
	switch conn.Dialector.Name() {
	case "postgres":
		myDialectFunction = postgres.TheFunctions()
	case "sqlite", "sqlite3":
		myDialectFunction = sqlite.TheFunctions()
	case "sqlserver":
		myDialectFunction = sqlserver.TheFunctions()
	case "bigquery":
		myDialectFunction = bigQuery.TheFunctions()
	default:
		err = fmt.Errorf("Dialect '%s' not fount", conn.Dialector.Name())
		return
	}

	ReservedWords = append(ReservedWords, myDialectFunction.ReservedWords...)
	return
}

func dataOutput(columns []structs.Data, rows *sql.Rows, columnsType []structs.ColumnType) ([]interface{}, error) {
	out := make([]interface{}, 0)

	for rows.Next() {
		fields, err := getFields(rows, columnsType)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})

		for j := 0; j < len(columns); j++ {
			column := columns[j]
			var dt string
			if column.Dt == nil {
				return nil, fmt.Errorf("Dt cannot be nil in column[%v]", j)
			}

			vType := reflect.TypeOf(column.Dt)
			if vType.String() == "string" {
				dt = column.Dt.(string)
			} else {
				dt = strconv.Itoa(column.Dt.(int))
			}

			db := strings.Replace(column.Db, myDialectFunction.EscapeChar, "", -1)

			if column.Formatter != nil {
				var err error
				row[dt], err = column.Formatter(fields[db], fields)
				if err != nil {
					return nil, err
				}
			} else {
				row[dt] = fields[db]
			}

		}
		out = append(out, row)
	}

	return out, nil
}

func drawNumber(c Controller) int {
	draw, err := strconv.Atoi(c.GetString("draw"))
	if err != nil {
		return 0
	}
	return draw
}

func Flated(whereArray []string) string {
	query := ""
	for _, where := range whereArray {
		if query != "" && where != "" {
			query += " AND "
		}
		query += where
	}
	return query
}

func buildSelectAndType(table string, join []structs.JoinData, conn *gorm.DB) (query string, fieldAlias map[string]string, err error) {
	query = fmt.Sprintf("%s.*", table)

	fieldAlias = make(map[string]string)

	columnsType, err := getSimpleBinding(conn, table)
	for _, columnInfo := range columnsType {
		fieldAlias[columnInfo.ColumnName] = myDialectFunction.EscapeChar + columnInfo.ColumnName + myDialectFunction.EscapeChar
	}
	if len(join) == 0 {
		return
	}
	subQuery, err := addFieldsSelect(table, table, conn, &fieldAlias)
	query += subQuery
	for _, tableData := range join {
		alias := tableData.Alias
		if alias == "" {
			alias = tableData.Table
		}
		subQuery, err = addFieldsSelect(tableData.Table, alias, conn, &fieldAlias)
		query += subQuery
	}

	return
}
func buildType(table string, conn *gorm.DB) (fieldAlias map[string]string) {

	fieldAlias = make(map[string]string)

	columnsType, _ := getSimpleBinding(conn, table)
	for _, columnInfo := range columnsType {
		fieldAlias[columnInfo.ColumnName] = myDialectFunction.EscapeChar + columnInfo.ColumnName + myDialectFunction.EscapeChar
	}

	addFieldsSelect(table, table, conn, &fieldAlias)

	return
}

func addFieldsSelect(table, alias string, conn *gorm.DB, fieldAlias *map[string]string) (query string, err error) {
	columnsType, err := getSimpleBinding(conn, table)
	for _, columnInfo := range columnsType {
		originalName := myDialectFunction.EscapeChar + alias + myDialectFunction.EscapeChar +
			"." + myDialectFunction.EscapeChar + columnInfo.ColumnName + myDialectFunction.EscapeChar

		aliasName := alias + myDialectFunction.AliasSeparator +
			columnInfo.ColumnName

		(*fieldAlias)[aliasName] = originalName

		query += ", " + originalName + " AS " + myDialectFunction.EscapeChar + aliasName + myDialectFunction.EscapeChar
	}
	return
}

func setJoins(joins []structs.JoinData) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, join := range joins {
			db = db.Joins(join.Query)
		}
		return db
	}
}

func setQuery(db *gorm.DB, query, logic string, param interface{}) *gorm.DB {
	if logic == "where" {
		if param == "" {
			db = db.Where(query)
		} else {
			db = db.Where(query, param)
		}
	} else {
		if param == "" {
			db = db.Or(query)
		} else {
			db = db.Or(query, param)
		}
	}

	return db
}

func setGlobalQuery(db *gorm.DB, query string, param interface{}, first bool) *gorm.DB {
	logic := "or"
	if first {
		logic = "where"
	}
	return setQuery(db, query, logic, param)
}

// database func
func filterGlobal(c Controller, columns []structs.Data, columnsType []structs.ColumnType, db *gorm.DB) *gorm.DB {

	str := c.GetString("search[value]")
	if str == "" {
		return db
	}

	requestRegex := ParamToBool(c, "search[regex]")
	//all columns filtering
	var i int
	first := true
	for i = 0; ; i++ {
		keyColumnsI := fmt.Sprintf("columns[%d][data]", i)

		keyColumnsData := c.GetString(keyColumnsI)
		if keyColumnsData == "" {
			break
		}
		columnIdx := Search(columns, keyColumnsData)

		requestColumnQuery := fmt.Sprintf("columns[%d][searchable]", i)
		requestColumn := c.GetString(requestColumnQuery)

		if columnIdx > -1 && requestColumn == "true" {

			query, param := bindingTypes(str, columnsType, columns[columnIdx], requestRegex)
			if query == "" {
				continue
			}
			db = setGlobalQuery(db, query, param, first)
			first = false

		} else {
			if columnIdx < 0 && requestColumn == "true" {
				fmt.Printf("(002) Do you forgot searchable: false in column %v ? or wrong column name in client side\n (client field data: must be same than server side DT: field)\n", keyColumnsData)
			}
		}
	}
	return db

}

func filterIndividual(c Controller, columns []structs.Data, columnsType []structs.ColumnType, db *gorm.DB) *gorm.DB {
	// Individual column filtering
	var i int
	for i = 0; ; i++ {
		keyColumnsI := fmt.Sprintf("columns[%d][data]", i)

		keyColumnsData := c.GetString(keyColumnsI)
		if keyColumnsData == "" {
			break
		}

		columnIdx := Search(columns, keyColumnsData)

		requestColumnQuery := fmt.Sprintf("columns[%d][searchable]", i)
		requestColumn := c.GetString(requestColumnQuery)

		requestColumnQuery = fmt.Sprintf("columns[%d][search][value]", i)
		str := c.GetString(requestColumnQuery)
		if columnIdx > -1 && requestColumn == "true" && str != "" {
			requestRegexQuery := fmt.Sprintf("columns[%d][search][regex]", i)
			requestRegex, err := strconv.ParseBool(c.GetString(requestRegexQuery))
			if err != nil {
				requestRegex = false
			}

			query, param := bindingTypes(str, columnsType, columns[columnIdx], requestRegex)
			if query == "" {
				continue
			}
			db = setQuery(db, query, "where", param)

		} else {
			if columnIdx < 0 && requestColumn == "true" {
				fmt.Printf("(001) Do you forgot searchable: false in column %v ? or wrong column name in client side\n (client field data: must be same than server side DT: field)\n", keyColumnsData)
			}
		}
	}
	return db

}

func order(c Controller, columns []structs.Data, columnsType []structs.ColumnType) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if c.GetString("order[0][column]") != "" {
			var i int
			for i = 0; ; i++ {
				columnIdxTittle := fmt.Sprintf("order[%d][column]", i)

				columnIdxOrder := c.GetString(columnIdxTittle)
				if columnIdxOrder == "" {
					break
				}

				columnIdxTittle = fmt.Sprintf("columns[%s][data]", columnIdxOrder)
				requestColumnData := c.GetString(columnIdxTittle)
				columnIdx := Search(columns, requestColumnData)

				columnIdxTittle = fmt.Sprintf("columns[%s][orderable]", columnIdxOrder)

				if columnIdx > -1 && c.GetString(columnIdxTittle) == "true" {

					column := columns[columnIdx]
					columnIdxTittle = fmt.Sprintf("order[%d][dir]", i)
					requestColumnData = c.GetString(columnIdxTittle)

					query := myDialectFunction.Order(column.Db, requestColumnData, columnsType)

					db = db.Order(query)
				} else {
					if columnIdx < 0 && c.GetString(columnIdxTittle) == "true" {
						fmt.Printf("(003) Do you forgot orderable: false in column %v ?\n", columnIdxOrder)
					}
				}
			}
		}
		return db
	}
}

func limit(c Controller) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		start, err := strconv.Atoi(c.GetString("start"))
		if err != nil || start < 0 {
			start = 0
		}

		length, err := strconv.Atoi(c.GetString("length"))

		if err != nil || length < 0 {
			length = 10
		}
		if length == 0 {
			length = 1
		}

		return db.Offset(start).Limit(length)
	}
}

func Search(column []structs.Data, keyColumnsI string) int {
	var i int
	for i = 0; i < len(column); i++ {
		data := column[i]
		if data.Dt == nil {
			continue
		}
		var field string
		vType := reflect.TypeOf(data.Dt)
		if vType.String() == "string" {
			field = data.Dt.(string)
		} else {
			field = strconv.Itoa(data.Dt.(int))
		}

		if field == keyColumnsI {
			return i
		}
	}
	return -1
}

// check if searchable field is string
func bindingTypes(value string, columnsType []structs.ColumnType, column structs.Data, isRegEx bool) (string, interface{}) {
	columndb := column.Db
	for _, columnInfo := range columnsType {
		if strings.Replace(columndb, myDialectFunction.EscapeChar, "", -1) == columnInfo.ColumnName {
			return myDialectFunction.BindingTypesQuery(columnInfo.Type,
				CheckReserved(columndb), value, columnInfo, isRegEx, column)
		}
	}

	return "", ""
}

func getFields(rows *sql.Rows, columnsType []structs.ColumnType) (map[string]interface{}, error) {

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	length := len(columns)
	current := makeResultReceiver(length)
	columnsTypeRows, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	err = rows.Scan(current...)
	if err != nil {
		return nil, err
	}

	value := make(map[string]interface{})
	for i := 0; i < length; i++ {
		key := columns[i]
		val := *(current[i]).(*interface{})
		if val == nil {
			value[key] = nil
			continue
		}
		vType := reflect.TypeOf(val)
		typeColum := findType(columnsTypeRows[i].Name(), columnsType)
		searching := columnsTypeRows[i].DatabaseTypeName()
		value[key], err = myDialectFunction.ParseData(searching, key, val, vType, typeColum)
		if err != nil {
			return nil, err
		}

	}
	return value, nil
}

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

func makeResultReceiver(length int) []interface{} {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var current interface{}
		current = struct{}{}
		result = append(result, &current)
	}
	return result
}

func getSimpleBinding(db *gorm.DB, table string) ([]structs.ColumnType, error) {

	rows, err := db.Select("*").
		Table(table).
		Limit(1).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnsType, err := rows.ColumnTypes()

	if err != nil {
		return nil, err
	}

	//build binding
	var types []structs.ColumnType
	for _, element := range columnsType {
		var ct structs.ColumnType
		ct.ColumnName = element.Name()

		types = append(types, ct)
	}

	return types, nil
}

func initBinding(db *gorm.DB, selectQuery, table string, whereJoin []structs.JoinData, fieldAlias map[string]string) ([]structs.ColumnType, error) {
	columnsType, err := defaultBinding(db, selectQuery, table, whereJoin)

	if err != nil {
		return nil, err
	}

	columnsTypesByDialect, err := dialectBinding(db, table, whereJoin)

	if err != nil {
		return nil, err
	}

	//build binding
	var types []structs.ColumnType
	for _, element := range columnsType {
		var ct structs.ColumnType
		ct.ColumnName = element.Name()
		ct.OriginalName = fieldAlias[element.Name()]

		customType, existCustomType := columnsTypesByDialect[element.Name()]
		if existCustomType {
			ct.Type = customType
		} else {
			ct.Type = element.DatabaseTypeName()
		}
		types = append(types, ct)
	}
	return types, nil
}

func defaultBinding(db *gorm.DB, selectQuery, table string, whereJoin []structs.JoinData) ([]*sql.ColumnType, error) {
	//get bind types
	rows, err := db.Select(selectQuery).
		Table(table).
		Scopes(
			setJoins(whereJoin),
		).
		Limit(1).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows.ColumnTypes()
}
func dialectBinding(db *gorm.DB, table string, whereJoin []structs.JoinData) (typeReturn map[string]string, err error) {
	dialectTypes := myDialectFunction.BindTypes(db, table)
	typeReturn = dialectTypes

	for val, dialectType := range dialectTypes {
		typeReturn[table+myDialectFunction.AliasSeparator+val] = dialectType
	}

	for _, join := range whereJoin {
		dialectTypes := myDialectFunction.BindTypes(db, join.Table)
		for val, dialectType := range dialectTypes {
			typeReturn[join.Table+myDialectFunction.AliasSeparator+val] = dialectType
		}
	}

	return dialectTypes, nil
}

// CheckReserved Skip reserved words
func CheckReserved(columnName string) string {
	if isReserved(columnName) {
		return myDialectFunction.ParseReservedField(columnName)
	}
	return columnName
}

func findType(columnName string, columnsType []structs.ColumnType) structs.ColumnType {

	for _, columnType := range columnsType {
		if columnType.ColumnName == columnName {
			return columnType
		}
	}
	return columnsType[0]
}
