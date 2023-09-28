package db_test

import (
	test "github.com/juaismar/go-gormssp/test"
	dbTest "github.com/juaismar/go-gormssp/test/dialects/bigQuery"
	"github.com/juaismar/go-gormssp/test/dialects/bigQuery/db"

	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
)

// call test for global test, or dbTest for custom test
var _ = Describe("Test BigQuery", func() {
	db := db.OpenDB()
	dbTest.ComplexFunctionTest(db)
	test.RegExpTest(db)
	dbTest.Types(db)
	dbTest.SimpleFunctionTest(db)
	test.Errors(db)
})
