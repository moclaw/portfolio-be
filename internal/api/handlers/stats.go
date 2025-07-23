package handlers

import (
	"net/http"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	projectService     services.ProjectService
	experienceService  services.ExperienceService
	technologyService  services.TechnologyService
	serviceService     services.ServiceService
	testimonialService services.TestimonialService
	contactService     *services.ContactService
}

func NewStatsHandler(
	projectService services.ProjectService,
	experienceService services.ExperienceService,
	technologyService services.TechnologyService,
	serviceService services.ServiceService,
	testimonialService services.TestimonialService,
	contactService *services.ContactService,
) *StatsHandler {
	return &StatsHandler{
		projectService:     projectService,
		experienceService:  experienceService,
		technologyService:  technologyService,
		serviceService:     serviceService,
		testimonialService: testimonialService,
		contactService:     contactService,
	}
}

type CountsResponse struct {
	Projects       int64 `json:"projects"`
	Experiences    int64 `json:"experiences"`
	Technologies   int64 `json:"technologies"`
	Services       int64 `json:"services"`
	Testimonials   int64 `json:"testimonials"`
	Contacts       int64 `json:"contacts"`
	UnreadContacts int64 `json:"unread_contacts"`
}

// GetCounts godoc
// @Summary Get counts for all entities
// @Description Get the total count of projects, experiences, technologies, services, testimonials, contacts, and unread contacts
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=CountsResponse}
// @Failure 500 {object} utils.Response
// @Router /api/stats/counts [get]
func (h *StatsHandler) GetCounts(c *gin.Context) {
	projectsCount, err := h.projectService.GetProjectsCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get projects count", err)
		return
	}

	experiencesCount, err := h.experienceService.GetExperiencesCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get experiences count", err)
		return
	}

	technologiesCount, err := h.technologyService.GetTechnologiesCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get technologies count", err)
		return
	}

	servicesCount, err := h.serviceService.GetServicesCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get services count", err)
		return
	}

	testimonialsCount, err := h.testimonialService.GetTestimonialsCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get testimonials count", err)
		return
	}

	contactsCount, err := h.contactService.GetContactsCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get contacts count", err)
		return
	}

	unreadContactsCount, err := h.contactService.GetUnreadCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get unread contacts count", err)
		return
	}

	countsResponse := CountsResponse{
		Projects:       projectsCount,
		Experiences:    experiencesCount,
		Technologies:   technologiesCount,
		Services:       servicesCount,
		Testimonials:   testimonialsCount,
		Contacts:       contactsCount,
		UnreadContacts: unreadContactsCount,
	}

	utils.SuccessResponse(c, "Counts retrieved successfully", countsResponse)
}
