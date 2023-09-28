package db_test

import (
	"github.com/juaismar/go-gormssp/test"

	"github.com/juaismar/go-gormssp/test/dialects/sqlite/db"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Test SQLite", func() {
	db := db.OpenDB()

	test.ComplexFunctionTest(db)
	//TODO uncoment when work
	//test.RegExpTest(db)
	test.Types(db)
	test.SimpleFunctionTest(db)
	//TODO test id "INTEGER" type
	test.Errors(db)
})
