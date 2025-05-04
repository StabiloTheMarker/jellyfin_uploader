package util

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jellyfin_uploader/models"
	"log"
	"os"
	"strings"
	"time"
)

var DB *gorm.DB

func ReadEnvFile() map[string]string {
	bytes, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	content := string(bytes)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	env := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		splitLine := strings.Split(line, "=")
		if len(splitLine) != 2 {
			log.Fatalf("Invalid line in .env file: %s", line)
		}
		env[splitLine[0]] = splitLine[1]
	}
	return env
}
func InitDB() {
	var err error
	env := ReadEnvFile()
	loc, _ := time.LoadLocation("Europe/Vienna") // or "America/New_York", etc.
	time.Local = loc
	dbDist := env["DB"]
	DB, err = gorm.Open(sqlite.Open(dbDist), &gorm.Config{})
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
