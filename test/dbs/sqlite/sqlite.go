package sqlite

import (
	databases "github.com/juaismar/go-gormssp/test/dbs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenDB return the Database connection
func OpenDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	databases.InitDB(db)

	return db
}
