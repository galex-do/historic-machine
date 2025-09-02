package main

import (
        "historical-events-backend/internal/config"
        "historical-events-backend/internal/database"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/handlers"
        "historical-events-backend/internal/services"
        "log"
        "net/http"
        "os"
)

func main() {
        // Load configuration
        cfg := config.Load()
        log.Println("Configuration loaded successfully")

        // Initialize database connection
        db, err := database.NewConnection(&cfg.Database)
        if err != nil {
                log.Fatal("Failed to connect to database:", err)
        }
        defer db.Close()

        // Initialize repositories
        eventRepo := repositories.NewEventRepository(db.DB)
        templateRepo := repositories.NewTemplateRepository(db.DB)
        tagRepo := repositories.NewTagRepository(db.DB)
        userRepo := repositories.NewUserRepository(db.DB)
        log.Println("Repositories initialized successfully")

        // Initialize services
        jwtSecret := os.Getenv("JWT_SECRET")
        if jwtSecret == "" {
                jwtSecret = "your-secret-key-change-in-production" // Default for development
        }
        authService := services.NewAuthService(userRepo, jwtSecret)
        log.Println("Services initialized successfully")

        // Initialize router with all handlers
        router := handlers.NewRouter(eventRepo, templateRepo, tagRepo, authService)
        
        // Setup routes
        httpHandler := router.SetupRoutes()
        
        // Start server
        serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
        log.Printf("Server starting on %s", serverAddr)
        log.Printf("API endpoints available at http://%s/api", serverAddr)
        
        if err := http.ListenAndServe(serverAddr, httpHandler); err != nil {
                log.Fatal("Server failed to start:", err)
        }
}