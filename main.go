package main

import (
	"github.com/gin-gonic/gin"

	"github.com/alrobwilloliver/animal-service-gin-dockertest/controller"
	"github.com/alrobwilloliver/animal-service-gin-dockertest/store"
)

func main() {
	store.InitDatabase()
	handler := controller.NewHandler(&store.Querier{})
	run(handler)
}

func run(h *controller.Handler) {
	router := gin.Default()
	router.GET("/animal", h.GetAnimals)
	router.POST("/animal", h.CreateAnimal)

	router.Run("localhost:8080")
}
