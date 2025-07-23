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

type ResourceHandler struct {
	service *services.ResourceService
}

func NewResourceHandler(service *services.ResourceService) *ResourceHandler {
	return &ResourceHandler{service: service}
}

// CreateResource godoc
// @Summary Create a new resource
// @Description Create a new resource with upload reference
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body models.ResourceCreateRequest true "Resource data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/resources [post]
func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var req models.ResourceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	resource, err := h.service.CreateResource(&req)
	if err != nil {
		if err.Error() == "upload not found" {
			utils.NotFoundResponse(c, "Upload not found")
		} else {
			utils.InternalErrorResponse(c, err)
		}
		return
	}

	utils.CreatedResponse(c, "Resource created successfully", resource)
}

// GetResource godoc
// @Summary Get resource by ID
// @Description Get a single resource record by its ID
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/resources/{id} [get]
func (h *ResourceHandler) GetResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid resource ID", err)
		return
	}

	resource, err := h.service.GetResourceByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Resource not found")
		return
	}

	utils.SuccessResponse(c, "Resource retrieved successfully", resource)
}

// GetAllResources godoc
// @Summary Get all resources
// @Description Get a list of all resource records
// @Tags resources
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param type query string false "Filter by resource type"
// @Param category query string false "Filter by category"
// @Param public query bool false "Filter public resources only"
// @Param search query string false "Search in name, description, or tags"
// @Success 200 {object} utils.PaginatedResponse
// @Failure 500 {object} utils.Response
// @Router /api/v1/resources [get]
func (h *ResourceHandler) GetAllResources(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	resourceType := c.Query("type")
	category := c.Query("category")
	publicOnly := c.Query("public") == "true"
	search := c.Query("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var resources []models.ResourceResponse

	// Handle different query types
	switch {
	case search != "":
		resources, err = h.service.SearchResources(search, limit, offset)
	case publicOnly:
		resources, err = h.service.GetPublicResources(limit, offset)
	case resourceType != "":
		resources, err = h.service.GetResourcesByType(models.ResourceType(resourceType), limit, offset)
	case category != "":
		resources, err = h.service.GetResourcesByCategory(category, limit, offset)
	default:
		resources, err = h.service.GetAllResources(limit, offset)
	}

	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	// For simplicity, calculate total pages based on returned results
	totalPages := int(math.Ceil(float64(len(resources)) / float64(limit)))
	if len(resources) == limit {
		totalPages = page + 1 // Assume there might be more pages
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int64(len(resources)),
		TotalPages: totalPages,
	}

	utils.PaginatedSuccessResponse(c, "Resources retrieved successfully", resources, pagination)
}

// UpdateResource godoc
// @Summary Update a resource
// @Description Update an existing resource record
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Resource ID"
// @Param resource body models.ResourceUpdateRequest true "Resource update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/resources/{id} [put]
func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid resource ID", err)
		return
	}

	var req models.ResourceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	resource, err := h.service.UpdateResource(uint(id), &req)
	if err != nil {
		if err.Error() == "resource not found" {
			utils.NotFoundResponse(c, "Resource not found")
		} else {
			utils.InternalErrorResponse(c, err)
		}
		return
	}

	utils.SuccessResponse(c, "Resource updated successfully", resource)
}

// DeleteResource godoc
// @Summary Delete a resource
// @Description Delete a resource record
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/resources/{id} [delete]
func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid resource ID", err)
		return
	}

	err = h.service.DeleteResource(uint(id))
	if err != nil {
		if err.Error() == "resource not found" {
			utils.NotFoundResponse(c, "Resource not found")
		} else {
			utils.InternalErrorResponse(c, err)
		}
		return
	}

	utils.SuccessResponse(c, "Resource deleted successfully", nil)
}

// DownloadResource godoc
// @Summary Download a resource
// @Description Download a resource file and increment download count
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/resources/{id}/download [post]
func (h *ResourceHandler) DownloadResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid resource ID", err)
		return
	}

	resource, err := h.service.GetResourceByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Resource not found")
		return
	}

	// Increment download count
	go h.service.IncrementDownloadCount(uint(id))

	// Return resource with download URL
	utils.SuccessResponse(c, "Resource download URL retrieved", map[string]interface{}{
		"resource":     resource,
		"download_url": resource.Upload.URL,
	})
}

// GetResourceStats godoc
// @Summary Get resource statistics
// @Description Get statistics about resources
// @Tags resources
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/resources/stats [get]
func (h *ResourceHandler) GetResourceStats(c *gin.Context) {
	stats, err := h.service.GetResourceStats()
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Resource statistics retrieved successfully", stats)
}

// RefreshExpiredURLs godoc
// @Summary Refresh expired URLs manually
// @Description Manually trigger the refresh of expired URLs
// @Tags resources
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/resources/refresh-urls [post]
func (h *ResourceHandler) RefreshExpiredURLs(c *gin.Context) {
	err := h.service.RefreshExpiredURLs()
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "URLs refreshed successfully", nil)
}
