package db_test

import (
	dbTest "github.com/juaismar/go-gormssp/test/dialects/bigQuery"
	"github.com/juaismar/go-gormssp/test/dialects/bigQuery/db"

	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
)

// call test for global test, or dbTest for custom test
var _ = Describe("Test BigQuery", func() {
	db := db.OpenDB()
	dbTest.ComplexFunctionTest(db)
	dbTest.RegExpTest(db)
	dbTest.Types(db)
	dbTest.SimpleFunctionTest(db)
	dbTest.Errors(db)
})
