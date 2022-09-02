package main

import (
	"github.com/gin-gonic/gin"

	"github.com/alrobwilloliver/animal-service-gin-dockertest/controller"
	"github.com/alrobwilloliver/animal-service-gin-dockertest/store"
)

func main() {
	store.InitDatabase()
	router := setUpRouter()
	router.Run(":8080")
}

func setUpRouter() *gin.Engine {
	h := controller.NewHandler(&store.Querier{})
	router := gin.Default()
	router.GET("/animal", h.GetAnimals)
	router.POST("/animal", h.CreateAnimal)

	return router
}
