package database

import (
	"database/sql"
	"fmt"
	"historical-events-backend/internal/config"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(cfg *config.DatabaseConfig) (*DB, error) {
	connectionString := getConnectionString(cfg)
	
	var db *sql.DB
	var err error
	
	// Retry connection with backoff for Docker/container environments
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Printf("Failed to open database connection: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		
		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		
		break
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to PostgreSQL database")
	return &DB{db}, nil
}

// getConnectionString builds the database connection string
func getConnectionString(cfg *config.DatabaseConfig) string {
	// Use DATABASE_URL if provided (for Replit/production)
	if cfg.DatabaseURL != "" {
		return cfg.DatabaseURL
	}
	
	// Build connection string from individual components
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Health checks if the database connection is healthy
func (db *DB) Health() error {
	return db.Ping()
}