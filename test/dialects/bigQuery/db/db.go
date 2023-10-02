package db

import (
	"time"

	"gorm.io/driver/bigquery"
	"gorm.io/gorm"
)

// OpenDB return the Database connection TODO
func OpenDB() *gorm.DB {
	db, err := gorm.Open(bigquery.Open("bigquery://testbigquerygorm/prueba"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//InitDB(db)

	return db
}

func InitDB(db *gorm.DB) {

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Pet{})

	db.Unscoped().Where("false_id >= 0").Delete(&User{})
	db.Unscoped().Where("false_id >= 0").Delete(&Pet{})

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

type CustomModel struct {
	FalseID   uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
