package main

import (
        "context"
        "historical-events-backend/internal/config"
        "historical-events-backend/internal/database"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/handlers"
        "historical-events-backend/internal/services"
        "log"
        "net/http"
        "os"
        "os/signal"
        "sync"
        "syscall"
        "time"
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
        datasetRepo := repositories.NewDatasetRepository(db.DB)
        supportRepo := repositories.NewSupportRepository(db.DB)
        log.Println("Repositories initialized successfully")

        // Initialize services
        jwtSecret := os.Getenv("JWT_SECRET")
        if jwtSecret == "" {
                jwtSecret = "your-secret-key-change-in-production" // Default for development
        }
        authService := services.NewAuthService(userRepo, jwtSecret)
        log.Println("Services initialized successfully")

        // Initialize metrics collector service
        metricsCollector := services.NewMetricsCollector(db.DB, log.New(os.Stdout, "[Metrics] ", log.LstdFlags))
        
        // Start metrics collection in background
        ctx, cancel := context.WithCancel(context.Background())
        var wg sync.WaitGroup
        wg.Add(1)
        go func() {
                defer wg.Done()
                metricsCollector.Start(ctx)
        }()

        // Initialize router with all handlers
        router := handlers.NewRouter(eventRepo, templateRepo, tagRepo, datasetRepo, authService, supportRepo)
        
        // Setup routes
        httpHandler := router.SetupRoutes()
        
        // Start HTTP server in a goroutine for graceful shutdown
        serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
        server := &http.Server{
                Addr:    serverAddr,
                Handler: httpHandler,
        }
        
        serverErrChan := make(chan error, 1)
        go func() {
                log.Printf("Server starting on %s", serverAddr)
                log.Printf("API endpoints available at http://%s/api", serverAddr)
                serverErrChan <- server.ListenAndServe()
        }()
        
        // Wait for interrupt signal for graceful shutdown
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
        
        select {
        case err := <-serverErrChan:
                if err != nil && err != http.ErrServerClosed {
                        log.Fatal("Server failed to start:", err)
                }
        case <-sigChan:
                log.Println("Shutting down gracefully...")
                
                // Cancel background services
                cancel()
                
                // Shutdown HTTP server with timeout
                shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
                defer shutdownCancel()
                
                if err := server.Shutdown(shutdownCtx); err != nil {
                        log.Printf("Server shutdown error: %v", err)
                }
                
                // Wait for background services to finish
                wg.Wait()
                log.Println("Server stopped")
        }
}