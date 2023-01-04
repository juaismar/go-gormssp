package ssp_test

import (
	"github.com/juaismar/go-gormssp/test/dbs/postgres"
	"github.com/juaismar/go-gormssp/test/dbs/sqlite"
	"github.com/juaismar/go-gormssp/test/dbs/sqlserver"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Test SQLITE", func() {
	db := sqlite.OpenDB()

	ComplexFunctionTest(db)
	//TODO uncoment when work
	//RegExpTest(db)
	Types(db)
	SimpleFunctionTest(db)
	//TODO test id "INTEGER" type
	Errors(db)
})

var _ = Describe("Test POSTGRES", func() {
	db := postgres.OpenDB()

	ComplexFunctionTest(db)
	RegExpTest(db)
	Types(db)
	SimpleFunctionTest(db)
	Errors(db)
})

var _ = Describe("Test SQLserver", func() {
	db := sqlserver.OpenDB()

	ComplexFunctionTest(db)
	//RegExpTest(db)
	Types(db)
	SimpleFunctionTest(db)
	Errors(db)
})

var _ = Describe("Test aux functions", func() {
	FunctionsTest()
})
