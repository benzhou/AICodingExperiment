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

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	dataSourceRepo := repository.NewDataSourceRepository()
	roleRepo := repository.NewRoleRepository()
	transactionRepo := repository.NewTransactionRepository()

	// Initialize services
	jwtService := services.NewJWTService()
	roleService := services.NewRoleService(roleRepo, userRepo)
	dataSourceService := services.NewDataSourceService(dataSourceRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	userService := services.NewUserService(userRepo, roleService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, roleService)
	dataSourceHandler := handlers.NewDataSourceHandler(dataSourceService, roleService)
	uploadHandler := handlers.NewUploadHandler(dataSourceService, transactionService, roleService)
	userHandler := handlers.NewUserHandler(userService, roleService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/", handlers.ServeStaticFiles)
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/api/v1/auth/token-info", authHandler.GetTokenInfo).Methods("GET")
	r.HandleFunc("/api/v1/auth/token-info-public", authHandler.GetTokenInfoPublic).Methods("GET")

	// Serve static files from the frontend build directory
	r.PathPrefix("/static/").HandlerFunc(handlers.ServeStaticFiles)

	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	protected.HandleFunc("/auth/google", authHandler.GoogleAuth).Methods("GET")
	protected.HandleFunc("/auth/google/callback", authHandler.GoogleCallback).Methods("GET")

	// Data source routes
	protected.HandleFunc("/datasources", dataSourceHandler.GetAllDataSources).Methods("GET")
	protected.HandleFunc("/datasources/search", dataSourceHandler.SearchDataSources).Methods("GET")
	protected.HandleFunc("/datasources", dataSourceHandler.CreateDataSource).Methods("POST")
	protected.HandleFunc("/datasources/{id}", dataSourceHandler.GetDataSourceByID).Methods("GET")
	protected.HandleFunc("/datasources/{id}", dataSourceHandler.UpdateDataSource).Methods("PUT")
	protected.HandleFunc("/datasources/{id}", dataSourceHandler.DeleteDataSource).Methods("DELETE")

	// Upload routes
	protected.HandleFunc("/uploads/transactions", uploadHandler.UploadTransactions).Methods("POST")

	// User management routes
	protected.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.GetUserById).Methods("GET")
	protected.HandleFunc("/users/{id}/roles", userHandler.GetUserRoles).Methods("GET")
	protected.HandleFunc("/users/{id}/roles", userHandler.UpdateUserRole).Methods("PUT")
	protected.HandleFunc("/users/{id}/admin", userHandler.SetAdminRole).Methods("PUT")
	protected.HandleFunc("/users", userHandler.CreateUserWithRole).Methods("POST")

	// Setup upload routes
	protected.HandleFunc("/uploads/preview", handlers.PreviewUploadHandler).Methods("POST")
	protected.HandleFunc("/uploads/process", handlers.ProcessUploadHandler).Methods("POST")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Accept"},
		AllowCredentials: true,
		MaxAge:           300,  // Maximum value not ignored by any of major browsers
		Debug:            true, // Enable debugging for CORS issues
	})

	// Start server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
