package handlers

import (
	"net/http"
	"strconv"

	"portfolio-be/internal/models"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TechnologyHandler struct {
	technologyService services.TechnologyService
}

func NewTechnologyHandler(technologyService services.TechnologyService) *TechnologyHandler {
	return &TechnologyHandler{
		technologyService: technologyService,
	}
}

// GetTechnologies
// @Summary Get all technologies
// @Description Get all active technologies ordered by order field
// @Tags technologies
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.TechnologyResponse}
// @Failure 500 {object} utils.Response
// @Router /api/technologies [get]
func (h *TechnologyHandler) GetTechnologies(c *gin.Context) {
	technologies, err := h.technologyService.GetActiveTechnologies()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get technologies", err)
		return
	}

	utils.SuccessResponse(c, "Technologies retrieved successfully", technologies)
}

// GetTechnology
// @Summary Get technology by ID
// @Description Get a specific technology by its ID
// @Tags technologies
// @Accept json
// @Produce json
// @Param id path int true "Technology ID"
// @Success 200 {object} utils.Response{data=models.TechnologyResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/technologies/{id} [get]
func (h *TechnologyHandler) GetTechnology(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid technology ID", err)
		return
	}

	technology, err := h.technologyService.GetTechnologyByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Technology not found", err)
		return
	}

	utils.SuccessResponse(c, "Technology retrieved successfully", technology)
}

// CreateTechnology
// @Summary Create new technology
// @Description Create a new technology record
// @Tags technologies
// @Accept json
// @Produce json
// @Param technology body models.TechnologyRequest true "Technology data"
// @Success 201 {object} utils.Response{data=models.TechnologyResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/technologies [post]
func (h *TechnologyHandler) CreateTechnology(c *gin.Context) {
	var req models.TechnologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	technology, err := h.technologyService.CreateTechnology(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create technology", err)
		return
	}

	utils.CreatedResponse(c, "Technology created successfully", technology)
}

// UpdateTechnology
// @Summary Update technology
// @Description Update an existing technology record
// @Tags technologies
// @Accept json
// @Produce json
// @Param id path int true "Technology ID"
// @Param technology body models.TechnologyRequest true "Technology data"
// @Success 200 {object} utils.Response{data=models.TechnologyResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/technologies/{id} [put]
func (h *TechnologyHandler) UpdateTechnology(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid technology ID", err)
		return
	}

	var req models.TechnologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	technology, err := h.technologyService.UpdateTechnology(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update technology", err)
		return
	}

	utils.SuccessResponse(c, "Technology updated successfully", technology)
}

// DeleteTechnology
// @Summary Delete technology
// @Description Soft delete a technology record
// @Tags technologies
// @Accept json
// @Produce json
// @Param id path int true "Technology ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/technologies/{id} [delete]
func (h *TechnologyHandler) DeleteTechnology(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid technology ID", err)
		return
	}

	err = h.technologyService.DeleteTechnology(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete technology", err)
		return
	}

	utils.SuccessResponse(c, "Technology deleted successfully", nil)
}
