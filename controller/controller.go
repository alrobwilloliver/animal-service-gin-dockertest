package controller

import (
	"github.com/alrobwilloliver/animal-service-gin-dockertest/model"
	"github.com/alrobwilloliver/animal-service-gin-dockertest/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	querier store.Queries
}

func NewHandler(q store.Queries) *Handler {
	return &Handler{
		querier: q,
	}
}

type CreateAnimalInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) GetAnimals(c *gin.Context) {
	animals, err := h.querier.GetAll(store.DB)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"message": "error getting animals",
			"error":   err.Error(),
		})
		return
	}
	c.IndentedJSON(200, gin.H{
		"message": "success",
		"data":    animals,
	})
}

func (h *Handler) CreateAnimal(c *gin.Context) {
	var input CreateAnimalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.IndentedJSON(400, gin.H{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}
	animal := model.Animal{
		Name: input.Name,
	}

	createdAnimal, err := h.querier.Create(store.DB, animal)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"message": "error creating animal",
			"error":   err.Error(),
		})
		return
	}
	c.IndentedJSON(201, gin.H{
		"message": "success",
		"data":    createdAnimal,
	})
}
