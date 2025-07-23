package handlers

import (
	"net/http"
	"portfolio-be/internal/models"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceService services.ServiceService
}

func NewServiceHandler(serviceService services.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceService: serviceService}
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service in the portfolio
// @Tags services
// @Accept json
// @Produce json
// @Param service body models.ServiceRequest true "Service data"
// @Success 201 {object} utils.Response{data=models.ServiceResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var request models.ServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	response, err := h.serviceService.CreateService(&request)
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.CreatedResponse(c, "Service created successfully", response)
}

// GetAllServices godoc
// @Summary Get all services
// @Description Get all services from the portfolio
// @Tags services
// @Produce json
// @Param active query boolean false "Filter by active status"
// @Success 200 {object} utils.Response{data=[]models.ServiceResponse}
// @Failure 500 {object} utils.Response
// @Router /api/v1/services [get]
func (h *ServiceHandler) GetAllServices(c *gin.Context) {
	activeOnly := c.Query("active")

	var responses []models.ServiceResponse
	var err error

	if activeOnly == "true" {
		responses, err = h.serviceService.GetActiveServices()
	} else {
		responses, err = h.serviceService.GetAllServices()
	}

	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Services retrieved successfully", responses)
}

// GetService godoc
// @Summary Get a service by ID
// @Description Get a specific service by its ID
// @Tags services
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} utils.Response{data=models.ServiceResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/services/{id} [get]
func (h *ServiceHandler) GetService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid service ID", err)
		return
	}

	response, err := h.serviceService.GetServiceByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Service not found")
		return
	}

	utils.SuccessResponse(c, "Service retrieved successfully", response)
}

// UpdateService godoc
// @Summary Update a service
// @Description Update an existing service by ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param service body models.ServiceRequest true "Updated service data"
// @Success 200 {object} utils.Response{data=models.ServiceResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid service ID", err)
		return
	}

	var request models.ServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.serviceService.UpdateService(uint(id), &request)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update service", err)
		return
	}

	utils.SuccessResponse(c, "Service updated successfully", response)
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service by ID
// @Tags services
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid service ID", err)
		return
	}

	err = h.serviceService.DeleteService(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete service", err)
		return
	}

	utils.SuccessResponse(c, "Service deleted successfully", nil)
}
