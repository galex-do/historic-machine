package services

import (
	"context"
	"database/sql"
	"historical-events-backend/pkg/metrics"
	"log"
	"time"
)

type MetricsCollector struct {
	db     *sql.DB
	logger *log.Logger
}

func NewMetricsCollector(db *sql.DB, logger *log.Logger) *MetricsCollector {
	return &MetricsCollector{
		db:     db,
		logger: logger,
	}
}

func (m *MetricsCollector) Start(ctx context.Context) {
	m.logger.Println("Starting metrics collector...")
	
	m.collect()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			m.logger.Println("Metrics collector stopped")
			return
		case <-ticker.C:
			m.collect()
		}
	}
}

func (m *MetricsCollector) collect() {
	m.collectEventCount()
	m.collectUserCount()
	m.collectTagCount()
	m.collectDatasetCount()
	m.collectTemplateCount()
	m.collectSessionCounts()
}

func (m *MetricsCollector) collectEventCount() {
	var count float64
	err := m.db.QueryRow("SELECT COUNT(*) FROM historical_events").Scan(&count)
	if err != nil {
		m.logger.Printf("Error collecting event count: %v", err)
		return
	}
	metrics.EventsTotal.Set(count)
}

func (m *MetricsCollector) collectUserCount() {
	var count float64
	err := m.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		m.logger.Printf("Error collecting user count: %v", err)
		return
	}
	metrics.UsersTotal.Set(count)
}

func (m *MetricsCollector) collectTagCount() {
	var count float64
	err := m.db.QueryRow("SELECT COUNT(*) FROM tags").Scan(&count)
	if err != nil {
		m.logger.Printf("Error collecting tag count: %v", err)
		return
	}
	metrics.TagsTotal.Set(count)
}

func (m *MetricsCollector) collectDatasetCount() {
	var count float64
	err := m.db.QueryRow("SELECT COUNT(*) FROM datasets").Scan(&count)
	if err != nil {
		m.logger.Printf("Error collecting dataset count: %v", err)
		return
	}
	metrics.DatasetsTotal.Set(count)
}

func (m *MetricsCollector) collectTemplateCount() {
	var count float64
	err := m.db.QueryRow("SELECT COUNT(*) FROM date_templates").Scan(&count)
	if err != nil {
		m.logger.Printf("Error collecting template count: %v", err)
		return
	}
	metrics.TemplatesTotal.Set(count)
}

func (m *MetricsCollector) collectSessionCounts() {
	activeWindow := time.Now().Add(-5 * time.Minute)
	
	var authenticatedCount float64
	err := m.db.QueryRow(`
		SELECT COUNT(*) FROM users 
		WHERE last_active_at IS NOT NULL 
		AND last_active_at > $1
	`, activeWindow).Scan(&authenticatedCount)
	if err != nil {
		m.logger.Printf("Error collecting authenticated session count: %v", err)
	} else {
		metrics.ActiveSessionsAuthenticated.Set(authenticatedCount)
	}
	
	var anonymousCount float64
	err = m.db.QueryRow(`
		SELECT COUNT(*) FROM anonymous_sessions 
		WHERE last_seen_at > $1
	`, activeWindow).Scan(&anonymousCount)
	if err != nil {
		m.logger.Printf("Error collecting anonymous session count: %v", err)
	} else {
		metrics.ActiveSessionsAnonymous.Set(anonymousCount)
	}
}
