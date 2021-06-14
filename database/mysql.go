package database

import (
	"fmt"
	"log"

	"github.com/release-trackers/gin/config"
	"github.com/release-trackers/gin/database/seed"
	"github.com/release-trackers/gin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitConnection() *gorm.DB {
	dsn := GetDbConnectionString()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error")
	}
	// log.Print("DB Connection successful")

	//For running project seeder
	//comment the code after seeding

	return db
}

func DbMigrate() {
	db := InitConnection()
	// Migrate the schema
	err := db.AutoMigrate(
		&models.Users{},
		&models.Release{},
		&models.Project{},
		&models.ReleaseProject{},
		&models.DLS{},
		&models.DlsProjects{},
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}
	log.Printf("Migration is successful!")

}

func DbSeed() {
	db := InitConnection()
	for _, seeds := range seed.All() {
		if err := seeds.Run(db); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seeds.Name, err)
		}
	}
	log.Printf("Seeding is successful!")
}

func GetDbConnectionString() string {
	conf := config.New().Database
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
	)
}
