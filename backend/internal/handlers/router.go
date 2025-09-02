package handlers

import (
        "historical-events-backend/internal/database/repositories"
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
}

// NewRouter creates a new router with all handlers
func NewRouter(eventRepo *repositories.EventRepository, templateRepo *repositories.TemplateRepository, tagRepo *repositories.TagRepository, authService *services.AuthService) *Router {
        return &Router{
                eventHandler:    NewEventHandler(eventRepo),
                templateHandler: NewTemplateHandler(templateRepo),
                tagHandler:      NewTagHandler(tagRepo),
                authHandler:     NewAuthHandler(authService),
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
        api.HandleFunc("/auth/login", router.authHandler.Login).Methods("POST")
        api.HandleFunc("/auth/register", router.authHandler.Register).Methods("POST")
        api.HandleFunc("/auth/logout", router.authHandler.Logout).Methods("POST")
        api.HandleFunc("/auth/me", router.authHandler.Me).Methods("GET")
        api.HandleFunc("/auth/change-password", router.authHandler.ChangePassword).Methods("POST")
        
        // Event routes (public read, auth required for create)
        api.HandleFunc("/events", router.authHandler.OptionalAuthMiddleware(router.eventHandler.GetAllEvents)).Methods("GET")
        api.HandleFunc("/events", router.authHandler.AuthMiddleware(router.eventHandler.CreateEvent)).Methods("POST")
        api.HandleFunc("/events/{id}", router.authHandler.OptionalAuthMiddleware(router.eventHandler.GetEventByID)).Methods("GET")
        
        // Spatial query routes
        api.HandleFunc("/events/bbox", router.eventHandler.GetEventsInBBox).Methods("GET")
        api.HandleFunc("/events/radius", router.eventHandler.GetEventsInRadius).Methods("GET")
        
        // Template routes
        api.HandleFunc("/date-template-groups", router.templateHandler.GetAllGroups).Methods("GET")
        api.HandleFunc("/date-templates/{group_id}", router.templateHandler.GetTemplatesByGroup).Methods("GET")
        api.HandleFunc("/date-templates", router.templateHandler.GetAllTemplates).Methods("GET")
        
        // Tag routes
        api.HandleFunc("/tags", router.tagHandler.GetAllTags).Methods("GET")
        api.HandleFunc("/tags", router.tagHandler.CreateTag).Methods("POST")
        api.HandleFunc("/tags/{id}", router.tagHandler.GetTagByID).Methods("GET")
        api.HandleFunc("/tags/{id}", router.tagHandler.UpdateTag).Methods("PUT")
        api.HandleFunc("/tags/{id}", router.tagHandler.DeleteTag).Methods("DELETE")
        
        // Event-Tag relationship routes
        api.HandleFunc("/events/{event_id}/tags/{tag_id}", router.tagHandler.AddTagToEvent).Methods("POST")
        api.HandleFunc("/events/{event_id}/tags/{tag_id}", router.tagHandler.RemoveTagFromEvent).Methods("DELETE")
        api.HandleFunc("/events/{event_id}/tags", router.tagHandler.SetEventTags).Methods("PUT")
        
        // Health check endpoint
        api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"status": "healthy"}`))
        }).Methods("GET")
        
        return r
}