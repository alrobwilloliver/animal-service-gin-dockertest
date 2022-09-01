package model

import "github.com/jinzhu/gorm"

type Animal struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Queries interface {
	Create(db *gorm.DB, animal Animal) (Animal, error)
	GetAll(db *gorm.DB) ([]Animal, error)
}

type Querier struct{}

func (q *Querier) Create(db *gorm.DB, animal Animal) (Animal, error) {
	if err := db.Create(&animal).Error; err != nil {
		return animal, err
	}
	return animal, nil
}

func (q *Querier) GetAll(db *gorm.DB) ([]Animal, error) {
	var animals []Animal
	if err := db.Find(&animals).Error; err != nil {
		return nil, err
	}
	return animals, nil
}
