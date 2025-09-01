package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Success sends a success response
func Success(w http.ResponseWriter, data interface{}, message ...string) {
	response := SuccessResponse{Data: data}
	if len(message) > 0 {
		response.Message = message[0]
	}
	JSON(w, http.StatusOK, response)
}

// Created sends a created response
func Created(w http.ResponseWriter, data interface{}, message ...string) {
	response := SuccessResponse{Data: data}
	if len(message) > 0 {
		response.Message = message[0]
	}
	JSON(w, http.StatusCreated, response)
}

// Error sends an error response
func Error(w http.ResponseWriter, statusCode int, message string, details ...string) {
	response := ErrorResponse{
		Error: message,
		Code:  statusCode,
	}
	if len(details) > 0 {
		response.Message = details[0]
	}
	JSON(w, statusCode, response)
}

// BadRequest sends a bad request response
func BadRequest(w http.ResponseWriter, message string, details ...string) {
	Error(w, http.StatusBadRequest, message, details...)
}

// InternalError sends an internal server error response
func InternalError(w http.ResponseWriter, message string, details ...string) {
	Error(w, http.StatusInternalServerError, message, details...)
}

// NotFound sends a not found response
func NotFound(w http.ResponseWriter, message string, details ...string) {
	Error(w, http.StatusNotFound, message, details...)
}