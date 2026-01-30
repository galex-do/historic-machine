package metrics

import (
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
        // HTTP metrics
        HttpRequestsTotal = promauto.NewCounterVec(
                prometheus.CounterOpts{
                        Name: "http_requests_total",
                        Help: "Total number of HTTP requests",
                },
                []string{"method", "path", "status"},
        )

        HttpRequestDuration = promauto.NewHistogramVec(
                prometheus.HistogramOpts{
                        Name:    "http_request_duration_seconds",
                        Help:    "HTTP request duration in seconds",
                        Buckets: prometheus.DefBuckets,
                },
                []string{"method", "path"},
        )

        HttpRequestsInFlight = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "http_requests_in_flight",
                        Help: "Number of HTTP requests currently being processed",
                },
        )

        // Event metrics
        EventsTotal = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "events_total",
                        Help: "Total number of events in database",
                },
        )

        EventsCreated = promauto.NewCounter(
                prometheus.CounterOpts{
                        Name: "events_created_total",
                        Help: "Total number of events created",
                },
        )

        EventsUpdated = promauto.NewCounter(
                prometheus.CounterOpts{
                        Name: "events_updated_total",
                        Help: "Total number of events updated",
                },
        )

        EventsDeleted = promauto.NewCounter(
                prometheus.CounterOpts{
                        Name: "events_deleted_total",
                        Help: "Total number of events deleted",
                },
        )

        // User and session metrics
        UsersTotal = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "users_total",
                        Help: "Total number of registered users",
                },
        )

        ActiveSessionsAuthenticated = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "active_sessions_authenticated",
                        Help: "Number of active authenticated user sessions",
                },
        )

        ActiveSessionsAnonymous = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "active_sessions_anonymous",
                        Help: "Number of active anonymous sessions",
                },
        )

        LoginAttemptsTotal = promauto.NewCounterVec(
                prometheus.CounterOpts{
                        Name: "login_attempts_total",
                        Help: "Total number of login attempts",
                },
                []string{"status"},
        )

        // Database metrics
        TagsTotal = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "tags_total",
                        Help: "Total number of tags in database",
                },
        )

        DatasetsTotal = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "datasets_total",
                        Help: "Total number of datasets in database",
                },
        )

        TemplatesTotal = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Name: "templates_total",
                        Help: "Total number of date templates in database",
                },
        )

        // API error metrics
        ApiErrorsTotal = promauto.NewCounterVec(
                prometheus.CounterOpts{
                        Name: "api_errors_total",
                        Help: "Total number of API errors",
                },
                []string{"type", "path"},
        )

        // Database query metrics
        DatabaseQueryDuration = promauto.NewHistogramVec(
                prometheus.HistogramOpts{
                        Name:    "database_query_duration_seconds",
                        Help:    "Database query duration in seconds",
                        Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0},
                },
                []string{"operation", "table"},
        )
)
