package test_test

import (
	"github.com/juaismar/go-gormssp/test"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test aux functions", Label("network"), func() {
	test.FunctionsTest()
})
