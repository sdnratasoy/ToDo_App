import (
	"log"
	"net/http"

	"todo-backend/database"
	"todo-backend/handlers"
	"todo-backend/middlewares"
)

func main() {

	database.Connect()
	http.HandleFunc("/api/v1/register", middlewares.CorsMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/api/v1/login", middlewares.CorsMiddleware(handlers.LoginHandler))
	http.HandleFunc("/api/v1/todos", middlewares.CorsMiddleware(middlewares.AuthMiddleware(handlers.TodosHandler)))
	http.HandleFunc("/api/v1/todos/", middlewares.CorsMiddleware(middlewares.AuthMiddleware(handlers.TodoByIdHandler)))

	log.Println("Server 8080 portunda calisiyor...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
