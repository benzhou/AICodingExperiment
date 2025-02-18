package main

import (
	"log"
	"net/http"
	"path/filepath"

	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env file from the project root
	if err := godotenv.Load(filepath.Join(".", ".env")); err != nil {
		log.Printf("Warning: .env file not found: %v\n", err)
	}

	// Initialize DB
	db.InitDB()

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository()
	authHandler := handlers.NewAuthHandler(userRepo)

	jwtService := services.NewJWTService()
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/", handlers.ServeStaticFiles)

	// Serve static files from the frontend build directory
	r.PathPrefix("/static/").HandlerFunc(handlers.ServeStaticFiles)

	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	protected.HandleFunc("/auth/token-info", authHandler.GetTokenInfo).Methods("GET")
	protected.HandleFunc("/auth/google", authHandler.GoogleAuth).Methods("GET")
	protected.HandleFunc("/auth/google/callback", authHandler.GoogleCallback).Methods("GET")

	// Existing routes
	protected.HandleFunc("/hello", handlers.HelloHandler).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		//AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Start server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
