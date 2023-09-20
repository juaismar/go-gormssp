package ssp_test

import (
	test "github.com/juaismar/go-gormssp/test"
	"github.com/juaismar/go-gormssp/test/dbs/postgres"
	"github.com/juaismar/go-gormssp/test/dbs/sqlite"
	"github.com/juaismar/go-gormssp/test/dbs/sqlserver"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Test SQLITE", func() {
	db := sqlite.OpenDB()

	test.ComplexFunctionTest(db)
	//TODO uncoment when work
	//test.RegExpTest(db)
	test.Types(db)
	test.SimpleFunctionTest(db)
	//TODO test id "INTEGER" type
	test.Errors(db)
})

var _ = Describe("Test POSTGRES", func() {
	db := postgres.OpenDB()

	test.ComplexFunctionTest(db)
	test.RegExpTest(db)
	test.Types(db)
	test.SimpleFunctionTest(db)
	test.Errors(db)
})

var _ = Describe("Test SQLserver", func() {
	db := sqlserver.OpenDB()

	test.ComplexFunctionTest(db)
	//test.RegExpTest(db)
	test.Types(db)
	test.SimpleFunctionTest(db)
	test.Errors(db)
})

var _ = Describe("Test aux functions", func() {
	test.FunctionsTest()
})
