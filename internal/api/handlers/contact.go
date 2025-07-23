package handlers

import (
	"net/http"
	"portfolio-be/internal/models"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	contactService *services.ContactService
}

func NewContactHandler(contactService *services.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// CreateContact godoc
// @Summary Create a new contact message
// @Description Submit a contact form message
// @Tags contact
// @Accept json
// @Produce json
// @Param contact body models.ContactRequest true "Contact message data"
// @Success 201 {object} utils.Response{data=models.ContactResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/contacts [post]
func (h *ContactHandler) CreateContact(c *gin.Context) {
	var req models.ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	contact, err := h.contactService.CreateContact(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contact", err)
		return
	}

	response := models.ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		Subject:   contact.Subject,
		Message:   contact.Message,
		Status:    contact.Status,
		IsActive:  contact.IsActive,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	utils.CreatedResponse(c, "Contact message sent successfully", response)
}

// GetContacts godoc
// @Summary Get all contact messages (Admin only)
// @Description Get all contact messages with pagination
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status" Enums(unread,read,replied)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.ContactResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/contacts [get]
func (h *ContactHandler) GetContacts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var contacts []models.ContactResponse
	var total int64
	var err error

	if status != "" {
		contacts, total, err = h.contactService.GetContactsByStatus(status, page, limit)
	} else {
		contacts, total, err = h.contactService.GetAllContacts(page, limit)
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get contacts", err)
		return
	}

	pagination := utils.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	c.JSON(http.StatusOK, utils.PaginatedResponse{
		Success:    true,
		Message:    "Contacts retrieved successfully",
		Data:       contacts,
		Pagination: pagination,
	})
}

// GetContact godoc
// @Summary Get a contact message by ID (Admin only)
// @Description Get a specific contact message by its ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Success 200 {object} utils.Response{data=models.ContactResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/contacts/{id} [get]
func (h *ContactHandler) GetContact(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	contact, err := h.contactService.GetContactByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Contact not found", err)
		return
	}

	// Mark as read when viewed
	h.contactService.MarkAsRead(uint(id))

	utils.SuccessResponse(c, "Contact retrieved successfully", contact)
}

// UpdateContact godoc
// @Summary Update a contact message (Admin only)
// @Description Update contact message information
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Param contact body models.ContactUpdateRequest true "Contact update data"
// @Success 200 {object} utils.Response{data=models.ContactResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/contacts/{id} [put]
func (h *ContactHandler) UpdateContact(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	var req models.ContactUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	contact, err := h.contactService.UpdateContact(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Contact not found", err)
		return
	}

	utils.SuccessResponse(c, "Contact updated successfully", contact)
}

// DeleteContact godoc
// @Summary Delete a contact message (Admin only)
// @Description Delete a contact message from the system
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/contacts/{id} [delete]
func (h *ContactHandler) DeleteContact(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	if err := h.contactService.DeleteContact(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Contact not found", err)
		return
	}

	utils.SuccessResponse(c, "Contact deleted successfully", nil)
}

// GetUnreadCount godoc
// @Summary Get unread contact count (Admin only)
// @Description Get the number of unread contact messages
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=map[string]int64}
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/contacts/unread-count [get]
func (h *ContactHandler) GetUnreadCount(c *gin.Context) {
	count, err := h.contactService.GetUnreadCount()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get unread count", err)
		return
	}

	utils.SuccessResponse(c, "Unread count retrieved successfully", map[string]int64{
		"unread_count": count,
	})
}

// MarkAsRead godoc
// @Summary Mark contact as read (Admin only)
// @Description Mark a contact message as read
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/contacts/{id}/mark-read [patch]
func (h *ContactHandler) MarkAsRead(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	if err := h.contactService.MarkAsRead(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Contact not found", err)
		return
	}

	utils.SuccessResponse(c, "Contact marked as read", nil)
}
