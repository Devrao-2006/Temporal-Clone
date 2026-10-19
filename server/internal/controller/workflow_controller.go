package controller

import (
	"encoding/json"
	"net/http"

	"temporal-clone/server/internal/models"
	"temporal-clone/server/internal/repository"
)

type WorkflowController struct {
	repo *repository.WorkflowRepo
}

func NewWorkflowController(repo *repository.WorkflowRepo) *WorkflowController {
	return &WorkflowController{repo: repo}
}

// Create handles (PUT) workflow/create
func (c *WorkflowController) Create(w http.ResponseWriter, r *http.Request) {
	var wf models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if wf.ID == "" || wf.Name == "" {
		WriteError(w, http.StatusBadRequest, "workflow 'id' and 'name' are required")
		return
	}
	if wf.Version == "" {
		wf.Version = "1.0"
	}
	if err := c.repo.Create(&wf); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]any{"status": "success", "data": wf})
}

// Update handles (PUT) workflow/update
func (c *WorkflowController) Update(w http.ResponseWriter, r *http.Request) {
	var wf models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if wf.ID == "" {
		WriteError(w, http.StatusBadRequest, "workflow 'id' is required")
		return
	}
	if err := c.repo.Update(&wf); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"status": "success", "data": wf})
}

// List handles (GET) workflow
func (c *WorkflowController) List(w http.ResponseWriter, r *http.Request) {
	workflows, err := c.repo.List()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if workflows == nil {
		workflows = []models.Workflow{}
	}
	WriteJSON(w, http.StatusOK, workflows)
}

// GetByID handles (GET) workflow/:id
func (c *WorkflowController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wf, err := c.repo.GetByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "workflow not found")
		return
	}
	WriteJSON(w, http.StatusOK, wf)
}
