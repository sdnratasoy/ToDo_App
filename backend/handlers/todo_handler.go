package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"todo-backend/database"
	"todo-backend/models"

	"github.com/google/uuid"
)


func TodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "GET":
		var todos []models.Todo
		database.DB.Order("created_at desc").Find(&todos)
		json.NewEncoder(w).Encode(todos)

	case "POST":
		var todo models.Todo
		json.NewDecoder(r.Body).Decode(&todo)

		if todo.Title == "" {
			http.Error(w, `{"message":"Title gerekli"}`, http.StatusBadRequest)
			return
		}

		todo.ID = uuid.New().String() 
		todo.CreatedAt = time.Now()   
		todo.Completed = false        

		database.DB.Create(&todo)

		json.NewEncoder(w).Encode(todo)

	default:
		http.Error(w, `{"message":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func TodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		http.Error(w, `{"message":"ID gerekli"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {

	case "GET":
		var todo models.Todo

		result := database.DB.First(&todo, "id = ?", id)

		if result.Error != nil {
			http.Error(w, `{"message":"Todo bulunamadi"}`, http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(todo)

	case "PUT":
		var todo models.Todo

		result := database.DB.First(&todo, "id = ?", id)
		if result.Error != nil {
			http.Error(w, `{"message":"Todo bulunamadi"}`, http.StatusNotFound)
			return
		}

		var updateData models.Todo
		json.NewDecoder(r.Body).Decode(&updateData)

		if updateData.Title != "" {
			todo.Title = updateData.Title
		}
		todo.Completed = updateData.Completed
		database.DB.Save(&todo)

		json.NewEncoder(w).Encode(todo)

	case "DELETE":
		result := database.DB.Delete(&models.Todo{}, "id = ?", id)

		if result.RowsAffected == 0 {
			http.Error(w, `{"message":"Todo bulunamadi"}`, http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(models.Response{Message: "Silindi"})

	default:
		http.Error(w, `{"message":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
