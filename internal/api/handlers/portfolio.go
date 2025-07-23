package handlers

import (
	"net/http"

	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	experienceService  services.ExperienceService
	serviceService     services.ServiceService
	technologyService  services.TechnologyService
	projectService     services.ProjectService
	testimonialService services.TestimonialService
}

func NewPortfolioHandler(
	experienceService services.ExperienceService,
	serviceService services.ServiceService,
	technologyService services.TechnologyService,
	projectService services.ProjectService,
	testimonialService services.TestimonialService,
) *PortfolioHandler {
	return &PortfolioHandler{
		experienceService:  experienceService,
		serviceService:     serviceService,
		technologyService:  technologyService,
		projectService:     projectService,
		testimonialService: testimonialService,
	}
}

// PortfolioData represents the complete portfolio data structure
type PortfolioData struct {
	Services     interface{} `json:"services"`
	Technologies interface{} `json:"technologies"`
	Experiences  interface{} `json:"experiences"`
	Testimonials interface{} `json:"testimonials"`
	Projects     interface{} `json:"projects"`
}

// GetPortfolioData
// @Summary Get complete portfolio data
// @Description Get all portfolio data in one request (services, technologies, experiences, testimonials, projects)
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=PortfolioData}
// @Failure 500 {object} utils.Response
// @Router /api/portfolio [get]
func (h *PortfolioHandler) GetPortfolioData(c *gin.Context) {
	// Get all data concurrently
	serviceChan := make(chan interface{}, 1)
	technologyChan := make(chan interface{}, 1)
	experienceChan := make(chan interface{}, 1)
	testimonialChan := make(chan interface{}, 1)
	projectChan := make(chan interface{}, 1)
	errorChan := make(chan error, 5)

	// Fetch services
	go func() {
		services, err := h.serviceService.GetActiveServices()
		if err != nil {
			errorChan <- err
			serviceChan <- nil
		} else {
			serviceChan <- services
		}
	}()

	// Fetch technologies
	go func() {
		technologies, err := h.technologyService.GetActiveTechnologies()
		if err != nil {
			errorChan <- err
			technologyChan <- nil
		} else {
			technologyChan <- technologies
		}
	}()

	// Fetch experiences
	go func() {
		experiences, err := h.experienceService.GetActiveExperiences()
		if err != nil {
			errorChan <- err
			experienceChan <- nil
		} else {
			experienceChan <- experiences
		}
	}()

	// Fetch testimonials
	go func() {
		testimonials, err := h.testimonialService.GetActiveTestimonials()
		if err != nil {
			errorChan <- err
			testimonialChan <- nil
		} else {
			testimonialChan <- testimonials
		}
	}()

	// Fetch projects
	go func() {
		projects, err := h.projectService.GetActiveProjects()
		if err != nil {
			errorChan <- err
			projectChan <- nil
		} else {
			projectChan <- projects
		}
	}()

	// Collect results
	portfolioData := &PortfolioData{}

	select {
	case err := <-errorChan:
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolio data", err)
		return
	case portfolioData.Services = <-serviceChan:
	}

	select {
	case err := <-errorChan:
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolio data", err)
		return
	case portfolioData.Technologies = <-technologyChan:
	}

	select {
	case err := <-errorChan:
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolio data", err)
		return
	case portfolioData.Experiences = <-experienceChan:
	}

	select {
	case err := <-errorChan:
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolio data", err)
		return
	case portfolioData.Testimonials = <-testimonialChan:
	}

	select {
	case err := <-errorChan:
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolio data", err)
		return
	case portfolioData.Projects = <-projectChan:
	}

	utils.SuccessResponse(c, "Portfolio data retrieved successfully", portfolioData)
}
