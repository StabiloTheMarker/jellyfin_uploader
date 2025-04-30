package util

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jellyfin_uploader/models"
	"time"
)

var DB *gorm.DB

func InitDB() {
	var err error
	loc, _ := time.LoadLocation("Europe/Vienna") // or "America/New_York", etc.
	time.Local = loc
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(&models.File{})
	if err != nil {
		panic("failed to migrate file")
	}

	err = DB.AutoMigrate(&models.UploadProcess{})
	if err != nil {
		panic("failed to migrate upload process")
	}
}
