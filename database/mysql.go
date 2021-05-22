package database

import (
	"fmt"
	"log"

	"github.com/release-trackers/gin/config"
	"github.com/release-trackers/gin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitConnection() *gorm.DB {
	dsn := getDbConnectionString()
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
	// for _, seeds := range seed.All() {
	// 	if err := seeds.Run(db); err != nil {
	// 		log.Fatalf("Running seed '%s', failed with error: %s", seeds.Name, err)
	// 	}
	// }

	return db
}

func getDbConnectionString() string {
	conf := config.New()
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DBName,
	)
}
