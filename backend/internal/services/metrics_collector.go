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

func (m *MetricsCollector) safeCount(query string, args ...interface{}) float64 {
        var count float64
        var err error
        if len(args) > 0 {
                err = m.db.QueryRow(query, args...).Scan(&count)
        } else {
                err = m.db.QueryRow(query).Scan(&count)
        }
        if err != nil {
                return 0
        }
        return count
}

func (m *MetricsCollector) collectEventCount() {
        metrics.EventsTotal.Set(m.safeCount("SELECT COUNT(*) FROM events"))
}

func (m *MetricsCollector) collectUserCount() {
        metrics.UsersTotal.Set(m.safeCount("SELECT COUNT(*) FROM users"))
}

func (m *MetricsCollector) collectTagCount() {
        metrics.TagsTotal.Set(m.safeCount("SELECT COUNT(*) FROM tags"))
}

func (m *MetricsCollector) collectDatasetCount() {
        metrics.DatasetsTotal.Set(m.safeCount("SELECT COUNT(*) FROM event_datasets"))
}

func (m *MetricsCollector) collectTemplateCount() {
        metrics.TemplatesTotal.Set(m.safeCount("SELECT COUNT(*) FROM date_templates"))
}

func (m *MetricsCollector) collectSessionCounts() {
        activeWindow := time.Now().Add(-5 * time.Minute)
        
        metrics.ActiveSessionsAuthenticated.Set(m.safeCount(`
                SELECT COUNT(*) FROM user_sessions 
                WHERE is_active = true 
                AND last_seen_at > $1
                AND expires_at > NOW()
        `, activeWindow))
        
        metrics.ActiveSessionsAnonymous.Set(m.safeCount(`
                SELECT COUNT(*) FROM anonymous_sessions 
                WHERE last_seen_at > $1
        `, activeWindow))
}

