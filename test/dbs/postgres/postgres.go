package postgres

import (
	//"github.com/jinzhu/gorm"

	databases "github.com/juaismar/go-gormssp/test/dbs"
	//"github.com/juaismar/go-gormssp/test/dbs/postgres"

	//"gorm.io/driver/mysql"
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
