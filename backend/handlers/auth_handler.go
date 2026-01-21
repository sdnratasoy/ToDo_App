package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"todo-backend/database"
	"todo-backend/middlewares"
	"todo-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"message":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"message":"Username ve password gerekli"}`, http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:       uuid.New().String(), 
		Username: req.Username,
		Password: req.Password,
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		http.Error(w, `{"message":"Kullanici zaten var"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(models.Response{Message: "Kayit basarili"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"message":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	var user models.User
	result := database.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user)

	if result.Error != nil {
		http.Error(w, `{"message":"Yanlis kullanici adi veya sifre"}`, http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), 
	})

	tokenString, err := token.SignedString(middlewares.JwtKey)
	if err != nil {
		http.Error(w, `{"message":"Token olusturulamadi"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.TokenResponse{Token: tokenString})
}
