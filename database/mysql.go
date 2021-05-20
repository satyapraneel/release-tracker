package database

import (
	"github.com/release-trackers/gin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)


func InitConnection() *gorm.DB{
	dsn := "root:secret@tcp(127.0.0.1:3307)/release_tracker?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error")
	}

	log.Print("DB Connection successful")
	// Migrate the schema
	db.AutoMigrate(
		&models.Users{},
		&models.Release{},
	)

	return db
}
