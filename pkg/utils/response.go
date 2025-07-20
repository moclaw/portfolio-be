package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:""`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success" example:"true"`
	Message    string      `json:"message" example:"Data retrieved successfully"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	TotalItems int64 `json:"total_items" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}

func PaginatedSuccessResponse(c *gin.Context, message string, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func ValidationErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, "Validation failed", err)
}

func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

func InternalErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusInternalServerError, "Internal server error", err)
}
