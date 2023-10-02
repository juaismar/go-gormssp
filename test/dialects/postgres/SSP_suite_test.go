package db_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSSP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SSP Postgres Suite")
}
