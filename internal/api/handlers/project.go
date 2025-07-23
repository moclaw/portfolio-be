package handlers

import (
	"net/http"
	"strconv"

	"portfolio-be/internal/models"
	"portfolio-be/internal/services"
	"portfolio-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetProjects
// @Summary Get all projects
// @Description Get all active projects ordered by order field
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.ProjectResponse}
// @Failure 500 {object} utils.Response
// @Router /api/projects [get]
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	projects, err := h.projectService.GetActiveProjects()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get projects", err)
		return
	}

	utils.SuccessResponse(c, "Projects retrieved successfully", projects)
}

// GetProject
// @Summary Get project by ID
// @Description Get a specific project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} utils.Response{data=models.ProjectResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	project, err := h.projectService.GetProjectByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Project not found", err)
		return
	}

	utils.SuccessResponse(c, "Project retrieved successfully", project)
}

// CreateProject
// @Summary Create new project
// @Description Create a new project record
// @Tags projects
// @Accept json
// @Produce json
// @Param project body models.ProjectRequest true "Project data"
// @Success 201 {object} utils.Response{data=models.ProjectResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req models.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	project, err := h.projectService.CreateProject(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create project", err)
		return
	}

	utils.CreatedResponse(c, "Project created successfully", project)
}

// UpdateProject
// @Summary Update project
// @Description Update an existing project record
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body models.ProjectRequest true "Project data"
// @Success 200 {object} utils.Response{data=models.ProjectResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	var req models.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	project, err := h.projectService.UpdateProject(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update project", err)
		return
	}

	utils.SuccessResponse(c, "Project updated successfully", project)
}

// DeleteProject
// @Summary Delete project
// @Description Soft delete a project record
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	err = h.projectService.DeleteProject(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete project", err)
		return
	}

	utils.SuccessResponse(c, "Project deleted successfully", nil)
}
