package db_test

import (
	"github.com/juaismar/go-gormssp/test"
	"github.com/juaismar/go-gormssp/test/dialects/sqlserver/db"

	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test SQLServer", func() {
	db := db.OpenDB()

	test.ComplexFunctionTest(db)
	//test.RegExpTest(db)
	test.Types(db)
	test.SimpleFunctionTest(db)
	test.Errors(db)
})
