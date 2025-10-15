package handlers

import (
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/models"
        "historical-events-backend/internal/services"
        "historical-events-backend/pkg/middleware"
        "net/http"

        "github.com/gorilla/mux"
)

// Router holds all the route handlers
type Router struct {
        eventHandler    *EventHandler
        templateHandler *TemplateHandler
        tagHandler      *TagHandler
        authHandler     *AuthHandler
        datasetHandler  *DatasetHandler
        tileHandler     *TileHandler
}

// NewRouter creates a new router with all handlers
func NewRouter(eventRepo *repositories.EventRepository, templateRepo *repositories.TemplateRepository, tagRepo *repositories.TagRepository, datasetRepo *repositories.DatasetRepository, authService *services.AuthService) *Router {
        return &Router{
                eventHandler:    NewEventHandler(eventRepo, tagRepo, datasetRepo),
                templateHandler: NewTemplateHandler(templateRepo),
                tagHandler:      NewTagHandler(tagRepo),
                authHandler:     NewAuthHandler(authService),
                datasetHandler:  NewDatasetHandler(datasetRepo, eventRepo),
                tileHandler:     NewTileHandler("./tile_cache"),
        }
}

// SetupRoutes configures all the routes and returns the router
func (router *Router) SetupRoutes() http.Handler {
        r := mux.NewRouter()
        
        // Add CORS middleware
        r.Use(middleware.CORS())
        
        // API routes
        api := r.PathPrefix("/api").Subrouter()
        
        // Authentication routes (public)
        api.HandleFunc("/auth/login", router.authHandler.Login).Methods("POST", "OPTIONS")
        api.HandleFunc("/auth/register", router.authHandler.Register).Methods("POST", "OPTIONS")
        api.HandleFunc("/auth/logout", router.authHandler.Logout).Methods("POST", "OPTIONS")
        api.HandleFunc("/auth/me", router.authHandler.AuthMiddleware(router.authHandler.Me)).Methods("GET", "OPTIONS")
        api.HandleFunc("/auth/change-password", router.authHandler.AuthMiddleware(router.authHandler.ChangePassword)).Methods("POST", "OPTIONS")
        
        // Event routes (public read, auth required for create/update/delete)
        api.HandleFunc("/events", router.authHandler.OptionalAuthMiddleware(router.eventHandler.GetAllEvents)).Methods("GET", "OPTIONS")
        api.HandleFunc("/events", router.authHandler.AuthMiddleware(router.eventHandler.CreateEvent)).Methods("POST", "OPTIONS")
        api.HandleFunc("/events/import", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.eventHandler.ImportEvents)).Methods("POST", "OPTIONS")
        api.HandleFunc("/events/{id}", router.authHandler.OptionalAuthMiddleware(router.eventHandler.GetEventByID)).Methods("GET", "OPTIONS")
        api.HandleFunc("/events/{id}", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.eventHandler.UpdateEvent)).Methods("PUT", "OPTIONS")
        api.HandleFunc("/events/{id}", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.eventHandler.DeleteEvent)).Methods("DELETE", "OPTIONS")
        
        // Spatial query routes
        api.HandleFunc("/events/bbox", router.eventHandler.GetEventsInBBox).Methods("GET", "OPTIONS")
        api.HandleFunc("/events/radius", router.eventHandler.GetEventsInRadius).Methods("GET", "OPTIONS")
        
        // Template routes
        api.HandleFunc("/date-template-groups", router.templateHandler.GetAllGroups).Methods("GET", "OPTIONS")
        api.HandleFunc("/date-templates/{group_id}", router.templateHandler.GetTemplatesByGroup).Methods("GET", "OPTIONS")
        api.HandleFunc("/date-templates", router.templateHandler.GetAllTemplates).Methods("GET", "OPTIONS")
        
        // Tag routes
        api.HandleFunc("/tags", router.tagHandler.GetAllTags).Methods("GET", "OPTIONS")
        api.HandleFunc("/tags", router.tagHandler.CreateTag).Methods("POST", "OPTIONS")
        api.HandleFunc("/tags/{id}", router.tagHandler.GetTagByID).Methods("GET", "OPTIONS")
        api.HandleFunc("/tags/{id}", router.tagHandler.UpdateTag).Methods("PUT", "OPTIONS")
        api.HandleFunc("/tags/{id}", router.tagHandler.DeleteTag).Methods("DELETE", "OPTIONS")
        
        // Dataset routes (admin only)
        api.HandleFunc("/datasets", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.datasetHandler.GetAllDatasets)).Methods("GET", "OPTIONS")
        api.HandleFunc("/datasets", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.datasetHandler.CreateDataset)).Methods("POST", "OPTIONS")
        api.HandleFunc("/datasets/{id}", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.datasetHandler.GetDatasetByID)).Methods("GET", "OPTIONS")
        api.HandleFunc("/datasets/{id}/export", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.datasetHandler.ExportDataset)).Methods("GET", "OPTIONS")
        api.HandleFunc("/datasets/{id}", router.authHandler.RequireAccessLevel(models.AccessLevelAdmin)(router.datasetHandler.DeleteDataset)).Methods("DELETE", "OPTIONS")
        
        // User management routes (super users only)
        api.HandleFunc("/users", router.authHandler.RequireAccessLevel(models.AccessLevelSuper)(router.authHandler.GetAllUsers)).Methods("GET", "OPTIONS")
        api.HandleFunc("/users", router.authHandler.RequireAccessLevel(models.AccessLevelSuper)(router.authHandler.CreateUser)).Methods("POST", "OPTIONS")
        api.HandleFunc("/users/{id}", router.authHandler.RequireAccessLevel(models.AccessLevelSuper)(router.authHandler.UpdateUser)).Methods("PUT", "OPTIONS")
        api.HandleFunc("/users/{id}", router.authHandler.RequireAccessLevel(models.AccessLevelSuper)(router.authHandler.DeleteUser)).Methods("DELETE", "OPTIONS")
        
        // Event-Tag relationship routes
        api.HandleFunc("/events/{event_id}/tags/{tag_id}", router.tagHandler.AddTagToEvent).Methods("POST", "OPTIONS")
        api.HandleFunc("/events/{event_id}/tags/{tag_id}", router.tagHandler.RemoveTagFromEvent).Methods("DELETE", "OPTIONS")
        api.HandleFunc("/events/{event_id}/tags", router.tagHandler.SetEventTags).Methods("PUT", "OPTIONS")
        
        // Tile proxy endpoint (public, no auth required)
        api.HandleFunc("/tiles", router.tileHandler.GetTile).Methods("GET", "OPTIONS")
        
        // Health check endpoint
        api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"status": "healthy"}`))
        }).Methods("GET", "OPTIONS")
        
        return r
}