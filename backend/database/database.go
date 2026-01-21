package database

import (
	"log"
	"os"
	"time"

	"todo-backend/models" 
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "sude")
	password := getEnv("DB_PASSWORD", "postgresude123")
	dbname := getEnv("DB_NAME", "tododb")
	port := getEnv("DB_PORT", "5432")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

	var err error

	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break 
		}
		log.Println("Veritabani baglantisi bekleniyor...")
		time.Sleep(3 * time.Second) 
	}

	if err != nil {
		log.Fatal("Veritabani baglantisi basarisiz:", err)
	}

	DB.AutoMigrate(&models.Todo{}, &models.User{})

	log.Println("Veritabani baglantisi basarili!")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
