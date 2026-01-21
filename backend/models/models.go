package models

import "time"

type Todo struct {
	ID        string    `json:"id" gorm:"primaryKey"`  
	Title     string    `json:"title"`                  
	Completed bool      `json:"completed"`              
	CreatedAt time.Time `json:"created_at"`             
}

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`   
	Username string `json:"username" gorm:"unique"` 
	Password string `json:"password"`               
}

type LoginRequest struct {
	Username string `json:"username"` 
	Password string `json:"password"` 
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Response struct {
	Message string `json:"message"`
}
