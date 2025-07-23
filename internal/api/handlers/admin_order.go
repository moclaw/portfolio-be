package handlers

import (
	"net/http"

	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AdminOrderHandler struct {
	projectService     services.ProjectService
	experienceService  services.ExperienceService
	technologyService  services.TechnologyService
	serviceService     services.ServiceService
	testimonialService services.TestimonialService
}

func NewAdminOrderHandler(
	projectService services.ProjectService,
	experienceService services.ExperienceService,
	technologyService services.TechnologyService,
	serviceService services.ServiceService,
	testimonialService services.TestimonialService,
) *AdminOrderHandler {
	return &AdminOrderHandler{
		projectService:     projectService,
		experienceService:  experienceService,
		technologyService:  technologyService,
		serviceService:     serviceService,
		testimonialService: testimonialService,
	}
}

type UpdateOrderRequest struct {
	Items []services.OrderItem `json:"items" binding:"required"`
}

// UpdateProjectsOrder
// @Summary Update projects order
// @Description Update the order of multiple projects
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateOrderRequest true "Order update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/projects/order [put]
func (h *AdminOrderHandler) UpdateProjectsOrder(c *gin.Context) {
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	err := h.projectService.UpdateProjectsOrder(req.Items)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update projects order", err)
		return
	}

	utils.SuccessResponse(c, "Projects order updated successfully", nil)
}

// UpdateExperiencesOrder
// @Summary Update experiences order
// @Description Update the order of multiple experiences
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateOrderRequest true "Order update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/experiences/order [put]
func (h *AdminOrderHandler) UpdateExperiencesOrder(c *gin.Context) {
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	err := h.experienceService.UpdateExperiencesOrder(req.Items)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update experiences order", err)
		return
	}

	utils.SuccessResponse(c, "Experiences order updated successfully", nil)
}

// UpdateTechnologiesOrder
// @Summary Update technologies order
// @Description Update the order of multiple technologies
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateOrderRequest true "Order update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/technologies/order [put]
func (h *AdminOrderHandler) UpdateTechnologiesOrder(c *gin.Context) {
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	err := h.technologyService.UpdateTechnologiesOrder(req.Items)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update technologies order", err)
		return
	}

	utils.SuccessResponse(c, "Technologies order updated successfully", nil)
}

// UpdateServicesOrder
// @Summary Update services order
// @Description Update the order of multiple services
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateOrderRequest true "Order update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/services/order [put]
func (h *AdminOrderHandler) UpdateServicesOrder(c *gin.Context) {
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	err := h.serviceService.UpdateServicesOrder(req.Items)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update services order", err)
		return
	}

	utils.SuccessResponse(c, "Services order updated successfully", nil)
}

// UpdateTestimonialsOrder
// @Summary Update testimonials order
// @Description Update the order of multiple testimonials
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateOrderRequest true "Order update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/testimonials/order [put]
func (h *AdminOrderHandler) UpdateTestimonialsOrder(c *gin.Context) {
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	err := h.testimonialService.UpdateTestimonialsOrder(req.Items)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update testimonials order", err)
		return
	}

	utils.SuccessResponse(c, "Testimonials order updated successfully", nil)
}
