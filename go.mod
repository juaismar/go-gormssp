module github.com/juaismar/go-gormssp

go 1.14

require (
	github.com/lib/pq v1.8.0
	//github.com/mattn/go-sqlite3 must be 1.14, explained in lib readme
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/nxadm/tail v1.4.5 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/satori/go.uuid v1.2.0
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gorm.io/driver/postgres v1.4.6
	gorm.io/driver/sqlite v1.4.4
	gorm.io/driver/sqlserver v1.4.1
	gorm.io/gorm v1.24.5
)
