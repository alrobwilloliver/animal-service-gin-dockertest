package store

import (
	"log"

	"github.com/alrobwilloliver/animal-service-gin/model"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver
)

var DB *gorm.DB

func InitDatabase() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(model.Animal{})
	DB = db
}
