package handlers

import (
	"math"
	"net/http"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *services.UploadService
}

func NewUploadHandler(service *services.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file to S3 and save record to database
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/uploads [post]
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// Parse multipart form with 32MB max memory
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse multipart form", err)
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No file provided", err)
		return
	}
	defer file.Close()

	// Check file size (max 10MB)
	if header.Size > 10*1024*1024 {
		utils.ErrorResponse(c, http.StatusBadRequest, "File size exceeds 10MB limit", nil)
		return
	}

	// Upload file
	upload, err := h.service.UploadFile(file, header)
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.CreatedResponse(c, "File uploaded successfully", upload)
}

// GetUpload godoc
// @Summary Get upload by ID
// @Description Get a single upload record by its ID
// @Tags upload
// @Accept json
// @Produce json
// @Param id path int true "Upload ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/uploads/{id} [get]
func (h *UploadHandler) GetUpload(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid upload ID", err)
		return
	}

	upload, err := h.service.GetUploadByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Upload not found")
		return
	}

	utils.SuccessResponse(c, "Upload retrieved successfully", upload)
}

// GetAllUploads godoc
// @Summary Get all uploads
// @Description Get a list of all upload records
// @Tags upload
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 500 {object} utils.Response
// @Router /api/v1/uploads [get]
func (h *UploadHandler) GetAllUploads(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	uploads, totalCount, err := h.service.GetUploadsWithCount(limit, offset)
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	// Calculate pagination properly
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}

	utils.PaginatedSuccessResponse(c, "Uploads retrieved successfully", uploads, pagination)
}

// GetAllUploadsWithSummary godoc
// @Summary Get all uploads with summary
// @Description Get a list of all upload records with summary statistics
// @Tags upload
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(12)
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/uploads/summary [get]
func (h *UploadHandler) GetAllUploadsWithSummary(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "12")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 12
	}

	offset := (page - 1) * limit

	result, err := h.service.GetAllUploadsWithSummary(limit, offset)
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	// Calculate pagination for the response
	totalItems := result.Summary.TotalFiles
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	// Create response with pagination
	response := map[string]interface{}{
		"data":       result.Uploads,
		"summary":    result.Summary,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, "Uploads retrieved successfully", response)
}

// DeleteUpload godoc
// @Summary Delete an upload
// @Description Delete an upload record and associated file from S3
// @Tags upload
// @Accept json
// @Produce json
// @Param id path int true "Upload ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/uploads/{id} [delete]
func (h *UploadHandler) DeleteUpload(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid upload ID", err)
		return
	}

	err = h.service.DeleteUpload(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			utils.NotFoundResponse(c, "Upload not found")
		} else {
			utils.InternalErrorResponse(c, err)
		}
		return
	}

	utils.SuccessResponse(c, "Upload deleted successfully", nil)
}
