package handlers

import (
	"net/http"
	"strconv"

	"portfolio-be/internal/models"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ExperienceHandler struct {
	experienceService services.ExperienceService
}

func NewExperienceHandler(experienceService services.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{
		experienceService: experienceService,
	}
}

// GetExperiences
// @Summary Get all experiences
// @Description Get all active experiences ordered by order field
// @Tags experiences
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.ExperienceResponse}
// @Failure 500 {object} utils.Response
// @Router /api/experiences [get]
func (h *ExperienceHandler) GetExperiences(c *gin.Context) {
	experiences, err := h.experienceService.GetActiveExperiences()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get experiences", err)
		return
	}

	utils.SuccessResponse(c, "Experiences retrieved successfully", experiences)
}

// GetExperience
// @Summary Get experience by ID
// @Description Get a specific experience by its ID
// @Tags experiences
// @Accept json
// @Produce json
// @Param id path int true "Experience ID"
// @Success 200 {object} utils.Response{data=models.ExperienceResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/experiences/{id} [get]
func (h *ExperienceHandler) GetExperience(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	experience, err := h.experienceService.GetExperienceByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Experience not found", err)
		return
	}

	utils.SuccessResponse(c, "Experience retrieved successfully", experience)
}

// CreateExperience
// @Summary Create new experience
// @Description Create a new experience record
// @Tags experiences
// @Accept json
// @Produce json
// @Param experience body models.ExperienceRequest true "Experience data"
// @Success 201 {object} utils.Response{data=models.ExperienceResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/experiences [post]
func (h *ExperienceHandler) CreateExperience(c *gin.Context) {
	var req models.ExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	experience, err := h.experienceService.CreateExperience(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create experience", err)
		return
	}

	utils.CreatedResponse(c, "Experience created successfully", experience)
}

// UpdateExperience
// @Summary Update experience
// @Description Update an existing experience record
// @Tags experiences
// @Accept json
// @Produce json
// @Param id path int true "Experience ID"
// @Param experience body models.ExperienceRequest true "Experience data"
// @Success 200 {object} utils.Response{data=models.ExperienceResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/experiences/{id} [put]
func (h *ExperienceHandler) UpdateExperience(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	var req models.ExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	experience, err := h.experienceService.UpdateExperience(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update experience", err)
		return
	}

	utils.SuccessResponse(c, "Experience updated successfully", experience)
}

// DeleteExperience
// @Summary Delete experience
// @Description Soft delete an experience record
// @Tags experiences
// @Accept json
// @Produce json
// @Param id path int true "Experience ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/experiences/{id} [delete]
func (h *ExperienceHandler) DeleteExperience(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid experience ID", err)
		return
	}

	err = h.experienceService.DeleteExperience(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete experience", err)
		return
	}

	utils.SuccessResponse(c, "Experience deleted successfully", nil)
}
