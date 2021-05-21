package database

import (
	"github.com/release-trackers/gin/database/seed"
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
		&models.Project{},
		&models.ReleaseProject{},
	)

	//For running project seeder
	//comment the code after seeding
	for _, seeds := range seed.All() {
		if err := seeds.Run(db); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seeds.Name, err)
		}
	}

	return db
}
