package store

import (
	"github.com/alrobwilloliver/animal-service-gin-dockertest/model"
	"github.com/jinzhu/gorm"
)

type Queries interface {
	Create(db *gorm.DB, animal model.Animal) (model.Animal, error)
	GetAll(db *gorm.DB) ([]model.Animal, error)
}

type Querier struct{}

func (q *Querier) Create(db *gorm.DB, animal model.Animal) (model.Animal, error) {
	if err := db.Create(&animal).Error; err != nil {
		return animal, err
	}
	return animal, nil
}

func (q *Querier) GetAll(db *gorm.DB) ([]model.Animal, error) {
	var animals []model.Animal
	if err := db.Find(&animals).Error; err != nil {
		return nil, err
	}
	return animals, nil
}
