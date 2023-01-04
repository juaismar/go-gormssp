package sqlserver

import (
	databases "github.com/juaismar/go-gormssp/test/dbs"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// OpenDB return the Database connection
func OpenDB() *gorm.DB {
	db, err := gorm.Open(sqlserver.Open("sqlserver://sqlserver:password@localhost:1433"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	databases.InitDB(db)

	return db
}
