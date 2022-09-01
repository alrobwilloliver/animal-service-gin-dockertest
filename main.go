package main

import (
	"github.com/gin-gonic/gin"

	"github.com/alrobwilloliver/animal-service-gin/controller"
	"github.com/alrobwilloliver/animal-service-gin/model"
	"github.com/alrobwilloliver/animal-service-gin/store"
)

func main() {

	store.InitDatabase()

	handler := controller.NewHandler(&model.Querier{})

	router := gin.Default()
	router.GET("/animal", handler.GetAnimals)
	router.POST("/animal", handler.CreateAnimal)

	router.Run("localhost:8080")
}
