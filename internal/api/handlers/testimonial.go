package handlers

import (
	"net/http"
	"strconv"

	"portfolio-be/internal/models"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TestimonialHandler struct {
	testimonialService services.TestimonialService
}

func NewTestimonialHandler(testimonialService services.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{
		testimonialService: testimonialService,
	}
}

// GetTestimonials
// @Summary Get all testimonials
// @Description Get all active testimonials ordered by order field
// @Tags testimonials
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.TestimonialResponse}
// @Failure 500 {object} utils.Response
// @Router /api/testimonials [get]
func (h *TestimonialHandler) GetTestimonials(c *gin.Context) {
	testimonials, err := h.testimonialService.GetActiveTestimonials()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get testimonials", err)
		return
	}

	utils.SuccessResponse(c, "Testimonials retrieved successfully", testimonials)
}

// GetTestimonial
// @Summary Get testimonial by ID
// @Description Get a specific testimonial by its ID
// @Tags testimonials
// @Accept json
// @Produce json
// @Param id path int true "Testimonial ID"
// @Success 200 {object} utils.Response{data=models.TestimonialResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/testimonials/{id} [get]
func (h *TestimonialHandler) GetTestimonial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid testimonial ID", err)
		return
	}

	testimonial, err := h.testimonialService.GetTestimonialByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Testimonial not found", err)
		return
	}

	utils.SuccessResponse(c, "Testimonial retrieved successfully", testimonial)
}

// CreateTestimonial
// @Summary Create new testimonial
// @Description Create a new testimonial record
// @Tags testimonials
// @Accept json
// @Produce json
// @Param testimonial body models.TestimonialRequest true "Testimonial data"
// @Success 201 {object} utils.Response{data=models.TestimonialResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/testimonials [post]
func (h *TestimonialHandler) CreateTestimonial(c *gin.Context) {
	var req models.TestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	testimonial, err := h.testimonialService.CreateTestimonial(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create testimonial", err)
		return
	}

	utils.CreatedResponse(c, "Testimonial created successfully", testimonial)
}

// UpdateTestimonial
// @Summary Update testimonial
// @Description Update an existing testimonial record
// @Tags testimonials
// @Accept json
// @Produce json
// @Param id path int true "Testimonial ID"
// @Param testimonial body models.TestimonialRequest true "Testimonial data"
// @Success 200 {object} utils.Response{data=models.TestimonialResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/testimonials/{id} [put]
func (h *TestimonialHandler) UpdateTestimonial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid testimonial ID", err)
		return
	}

	var req models.TestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	testimonial, err := h.testimonialService.UpdateTestimonial(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update testimonial", err)
		return
	}

	utils.SuccessResponse(c, "Testimonial updated successfully", testimonial)
}

// DeleteTestimonial
// @Summary Delete testimonial
// @Description Soft delete a testimonial record
// @Tags testimonials
// @Accept json
// @Produce json
// @Param id path int true "Testimonial ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/testimonials/{id} [delete]
func (h *TestimonialHandler) DeleteTestimonial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid testimonial ID", err)
		return
	}

	err = h.testimonialService.DeleteTestimonial(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete testimonial", err)
		return
	}

	utils.SuccessResponse(c, "Testimonial deleted successfully", nil)
}
