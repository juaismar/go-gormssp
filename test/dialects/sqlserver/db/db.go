package db

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// OpenDB return the Database connection
func OpenDB() *gorm.DB {
	db, err := gorm.Open(sqlserver.Open("sqlserver://sqlserver:password@localhost:1433"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	InitDB(db)

	return db
}

func InitDB(db *gorm.DB) {

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Pet{})

	db.Unscoped().Where("id >= 0").Delete(&User{})
	db.Unscoped().Where("id >= 0").Delete(&Pet{})

	fillData(db)
}

func fillData(db *gorm.DB) {

	for _, user := range GetDefaultUser() {
		db.Create(&user)
	}

	for _, pet := range GetDefaultPet() {
		db.Create(&pet)
	}

}
