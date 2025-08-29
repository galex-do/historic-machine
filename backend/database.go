package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func get_database_config() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "historical_events"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func connect_database() (*sql.DB, error) {
	config := get_database_config()
	
	connection_string := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Try to connect with retries for Docker environment
	var db *sql.DB
	var err error
	
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connection_string)
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

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func get_events_from_db(db *sql.DB) ([]HistoricalEvent, error) {
	query := `SELECT id, name, description, latitude, longitude, event_date, lens_type FROM events ORDER BY event_date DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var events []HistoricalEvent
	
	for rows.Next() {
		var event HistoricalEvent
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Latitude, &event.Longitude, &event.EventDate, &event.LensType)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	
	return events, nil
}

func create_event_in_db(db *sql.DB, event HistoricalEvent) (*HistoricalEvent, error) {
	query := `INSERT INTO events (name, description, latitude, longitude, event_date, lens_type) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	
	var created_event HistoricalEvent = event
	var created_at time.Time
	
	err := db.QueryRow(query, event.Name, event.Description, event.Latitude, event.Longitude, event.EventDate, event.LensType).
		Scan(&created_event.ID, &created_at)
	
	if err != nil {
		return nil, err
	}
	
	return &created_event, nil
}