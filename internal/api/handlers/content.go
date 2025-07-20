package handlers

import (
	"math"
	"net/http"
	"portfolio-be/internal/models"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	service *services.ContentService
}

func NewContentHandler(service *services.ContentService) *ContentHandler {
	return &ContentHandler{service: service}
}

// CreateContent godoc
// @Summary Create a new content
// @Description Create a new content item
// @Tags content
// @Accept json
// @Produce json
// @Param content body models.ContentRequest true "Content data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/contents [post]
func (h *ContentHandler) CreateContent(c *gin.Context) {
	var req models.ContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	content, err := h.service.CreateContent(req)
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.CreatedResponse(c, "Content created successfully", content)
}

// GetContent godoc
// @Summary Get a content by ID
// @Description Get a single content item by its ID
// @Tags content
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/contents/{id} [get]
func (h *ContentHandler) GetContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err)
		return
	}

	content, err := h.service.GetContentByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Content not found")
		return
	}

	utils.SuccessResponse(c, "Content retrieved successfully", content)
}

// GetAllContents godoc
// @Summary Get all contents
// @Description Get a list of all content items with optional filtering
// @Tags content
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query string false "Filter by category"
// @Param status query string false "Filter by status"
// @Success 200 {object} utils.PaginatedResponse
// @Failure 500 {object} utils.Response
// @Router /api/v1/contents [get]
func (h *ContentHandler) GetAllContents(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	category := c.Query("category")
	status := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	contents, err := h.service.GetAllContent(limit, offset, category, status)
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	// Get total count for pagination
	totalCount, err := h.service.GetContentCount()
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}

	utils.PaginatedSuccessResponse(c, "Contents retrieved successfully", contents, pagination)
}

// UpdateContent godoc
// @Summary Update a content
// @Description Update an existing content item
// @Tags content
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Param content body models.ContentRequest true "Updated content data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/contents/{id} [put]
func (h *ContentHandler) UpdateContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err)
		return
	}

	var req models.ContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	content, err := h.service.UpdateContent(uint(id), req)
	if err != nil {
		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Content not found")
		} else {
			utils.InternalErrorResponse(c, err)
		}
		return
	}

	utils.SuccessResponse(c, "Content updated successfully", content)
}

// DeleteContent godoc
// @Summary Delete a content
// @Description Delete an existing content item
// @Tags content
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/contents/{id} [delete]
func (h *ContentHandler) DeleteContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err)
		return
	}

	err = h.service.DeleteContent(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Content not found")
		} else {
			utils.InternalErrorResponse(c, err)
		}
		return
	}

	utils.SuccessResponse(c, "Content deleted successfully", nil)
}