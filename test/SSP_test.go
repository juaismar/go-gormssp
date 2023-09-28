package test_test

import (
	"github.com/juaismar/go-gormssp/test"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Test aux functions", func() {
	test.FunctionsTest()
})
