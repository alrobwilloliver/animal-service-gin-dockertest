package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/alrobwilloliver/animal-service-gin/controller"
	"github.com/alrobwilloliver/animal-service-gin/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type fakeQuerier struct {
	animal  model.Animal
	animals []model.Animal
	err     error
}

func (f *fakeQuerier) Create(db *gorm.DB, animal model.Animal) (model.Animal, error) {
	return f.animal, f.err
}

func (f *fakeQuerier) GetAll(db *gorm.DB) ([]model.Animal, error) {
	return f.animals, f.err
}

var jsonResponse map[string]interface{}

func TestCreateAnimal(t *testing.T) {
	t.Run("should successfully create dog", func(t *testing.T) {
		router := gin.Default()
		fakeQuerier := &fakeQuerier{}
		fakeQuerier.animal = model.Animal{
			ID:   1,
			Name: "dog",
		}
		handler := controller.NewHandler(fakeQuerier)
		router.POST("/animal", handler.CreateAnimal)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/animal", bytes.NewBuffer([]byte(`{"name": "dog"}`)))
		router.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Errorf("Expected status code %d, got %d", 201, w.Code)
		}
		if w.Body == nil {
			t.Error("Expected body, got nil")
		}
		if err := json.Unmarshal(w.Body.Bytes(), &jsonResponse); err != nil {
			t.Errorf("Expected json, got error: %s", err)
		}
		if jsonResponse["message"].(string) != "success" {
			t.Errorf("Expected message %s, got %s", "success", jsonResponse["message"].(string))
		}
		if jsonResponse["data"] == nil {
			t.Error("Expected data, got nil")
		}
		if jsonResponse["data"].(map[string]interface{})["name"] != "dog" {
			t.Errorf("Expected name %s, got %s", "dog", jsonResponse["data"].(map[string]interface{})["name"])
		}
		if jsonResponse["data"].(map[string]interface{})["id"] != float64(1) {
			t.Errorf("Expected id %f, got %d", float64(1), jsonResponse["data"].(map[string]interface{})["id"])
		}
	})
	t.Run("should return error when animal name is empty", func(t *testing.T) {
		router := gin.Default()
		fakeQuerier := &fakeQuerier{}
		handler := controller.NewHandler(fakeQuerier)
		router.POST("/animal", handler.CreateAnimal)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/animal", bytes.NewBuffer([]byte(`{"name": ""}`)))
		router.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Errorf("Expected status code %d, got %d", 400, w.Code)
		}
		if w.Body == nil {
			t.Error("Expected body, got nil")
		}
		if err := json.Unmarshal(w.Body.Bytes(), &jsonResponse); err != nil {
			t.Errorf("Expected json, got error: %s", err)
		}
		if jsonResponse["message"].(string) != "invalid input" {
			t.Errorf("Expected message %s, got %s", "invalid input", jsonResponse["message"].(string))
		}
		if jsonResponse["error"] != "Key: 'CreateAnimalInput.Name' Error:Field validation for 'Name' failed on the 'required' tag" {
			t.Errorf("Expected %s, got %s", "Key: 'CreateAnimalInput.Name' Error:Field validation for 'Name' failed on the 'required' tag", jsonResponse["error"])
		}
	})
	t.Run("should error when Create fails", func(t *testing.T) {
		router := gin.Default()
		fakeQuerier := &fakeQuerier{}
		fakeQuerier.err = errors.New("error")
		handler := controller.NewHandler(fakeQuerier)
		router.POST("/animal", handler.CreateAnimal)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/animal", bytes.NewBuffer([]byte(`{"name": "dog"}`)))
		router.ServeHTTP(w, req)

		if w.Code != 500 {
			t.Errorf("Expected status code %d, got %d", 500, w.Code)
		}
		if w.Body == nil {
			t.Error("Expected body, got nil")
		}

		if err := json.Unmarshal(w.Body.Bytes(), &jsonResponse); err != nil {
			t.Errorf("Expected json, got error: %s", err)
		}
		if jsonResponse["message"].(string) != "error creating animal" {
			t.Errorf("Expected message %s, got %s", "error creating animal", jsonResponse["message"].(string))
		}
	})
}

func TestGetAllAnimals(t *testing.T) {
	t.Run("should fail when GetAll returns an error", func(t *testing.T) {
		router := gin.Default()
		fakeQuerier := &fakeQuerier{}
		fakeQuerier.err = errors.New("error")
		handler := controller.NewHandler(fakeQuerier)
		router.GET("/animal", handler.GetAnimals)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/animal", nil)
		router.ServeHTTP(w, req)

		if w.Code != 500 {
			t.Errorf("Expected status code %d, got %d", 500, w.Code)
		}
		if w.Body == nil {
			t.Error("Expected body, got nil")
		}

		if err := json.Unmarshal(w.Body.Bytes(), &jsonResponse); err != nil {
			t.Errorf("Expected json, got error: %s", err)
		}
		if jsonResponse["message"].(string) != "error getting animals" {
			t.Errorf("Expected message: %s, got %s", "error getting animals", jsonResponse["message"].(string))
		}
		if jsonResponse["error"] != "error" {
			t.Errorf("Expected error: %s, got %s", "error", jsonResponse["error"])
		}
	})
	t.Run("should successfully get all animals", func(t *testing.T) {
		router := gin.Default()
		fakeQuerier := &fakeQuerier{}
		fakeQuerier.animals = []model.Animal{
			{
				ID:   1,
				Name: "dog",
			},
			{
				ID:   2,
				Name: "cat",
			},
		}
		handler := controller.NewHandler(fakeQuerier)
		router.GET("/animals", handler.GetAnimals)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/animals", nil)
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status code %d, got %d", 200, w.Code)
		}
		if w.Body == nil {
			t.Error("Expected body, got nil")
		}
		if err := json.Unmarshal(w.Body.Bytes(), &jsonResponse); err != nil {
			t.Errorf("Expected json, got error: %s", err)
		}
		if jsonResponse["message"].(string) != "success" {
			t.Errorf("Expected message %s, got %s", "success", jsonResponse["message"].(string))
		}
		if jsonResponse["data"].([]interface{})[0].(map[string]interface{})["name"] != "dog" {
			t.Errorf("Expected name %s, got %s", "dog", jsonResponse["data"].([]interface{})[0].(map[string]interface{})["name"])
		}
		if jsonResponse["data"].([]interface{})[0].(map[string]interface{})["id"] != float64(1) {
			t.Errorf("Expected id %f, got %d", float64(1), jsonResponse["data"].([]interface{})[0].(map[string]interface{})["id"])
		}
		if jsonResponse["data"].([]interface{})[1].(map[string]interface{})["name"] != "cat" {
			t.Errorf("Expected name %s, got %s", "cat", jsonResponse["data"].([]interface{})[1].(map[string]interface{})["name"])
		}
		if jsonResponse["data"].([]interface{})[1].(map[string]interface{})["id"] != float64(2) {
			t.Errorf("Expected id %f, got %d", float64(2), jsonResponse["data"].([]interface{})[1].(map[string]interface{})["id"])
		}
		if jsonResponse["message"].(string) != "success" {
			t.Errorf("Expected message %s, got %s", "success", jsonResponse["message"].(string))
		}
	})
}
