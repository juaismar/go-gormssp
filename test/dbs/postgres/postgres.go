package postgres

import (
	databases "github.com/juaismar/go-gormssp/test/dbs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OpenDB return the Database connection
func OpenDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	databases.InitDB(db)

	return db
}
