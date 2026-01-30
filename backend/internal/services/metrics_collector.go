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
        m.collectSessionDurations()
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

func (m *MetricsCollector) collectSessionDurations() {
        // Only observe sessions that became inactive recently (last 5-10 minutes, no recent heartbeat)
        // This avoids double-counting sessions on every collection cycle
        inactiveWindow := time.Now().Add(-10 * time.Minute)
        activeWindow := time.Now().Add(-5 * time.Minute)

        // Anonymous sessions that ended (no heartbeat in 5 min, but had activity in last 10 min)
        rows, err := m.db.Query(`
                SELECT EXTRACT(EPOCH FROM (last_seen_at - created_at)) / 60 as duration_minutes
                FROM anonymous_sessions
                WHERE last_seen_at IS NOT NULL
                AND last_seen_at < $1
                AND last_seen_at > $2
        `, activeWindow, inactiveWindow)
        if err != nil {
                return
        }
        defer rows.Close()

        for rows.Next() {
                var duration float64
                if err := rows.Scan(&duration); err == nil && duration > 0 {
                        metrics.SessionDurationMinutes.WithLabelValues("anonymous").Observe(duration)
                }
        }

        // Authenticated sessions that ended recently
        authRows, err := m.db.Query(`
                SELECT EXTRACT(EPOCH FROM (COALESCE(ended_at, last_seen_at) - created_at)) / 60 as duration_minutes
                FROM user_sessions
                WHERE ended_at IS NOT NULL
                AND ended_at > $1
        `, inactiveWindow)
        if err != nil {
                return
        }
        defer authRows.Close()

        for authRows.Next() {
                var duration float64
                if err := authRows.Scan(&duration); err == nil && duration > 0 {
                        metrics.SessionDurationMinutes.WithLabelValues("authenticated").Observe(duration)
                }
        }
}
