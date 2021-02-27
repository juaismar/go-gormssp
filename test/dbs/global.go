package databases

import (
	model "github.com/juaismar/go-gormssp/test/models"
	"gorm.io/gorm"
)

// InitDB clear and populate database
func InitDB(db *gorm.DB) {

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Pet{})

	db.Unscoped().Where("id >= 0").Delete(&model.User{})
	db.Unscoped().Where("id >= 0").Delete(&model.Pet{})
	fillData(db)
}

func fillData(db *gorm.DB) {

	for _, user := range model.GetDefaultUser() {
		db.Create(&user)
	}
	for _, pet := range model.GetDefaultPet() {
		db.Create(&pet)
	}

}
