package store

import (
	"log"

	"github.com/alrobwilloliver/animal-service-gin/model"
	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func InitDatabase() {
	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(model.Animal{})
	DB = db
}
