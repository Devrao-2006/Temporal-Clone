package controller

import (
	"encoding/json"
	"net/http"

	"temporal-clone/server/internal/models"
	"temporal-clone/server/internal/repository"
	"temporal-clone/server/internal/service"
)

type ExecutionController struct {
	repo    *repository.ExecutionRepo
	service *service.ExecutionService
}

func NewExecutionController(repo *repository.ExecutionRepo, svc *service.ExecutionService) *ExecutionController {
	return &ExecutionController{
		repo:    repo,
		service: svc,
	}
}

// Create handles (PUT) execution/create
func (c *ExecutionController) Create(w http.ResponseWriter, r *http.Request) {
	var req service.CreateExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.WorkflowID == "" {
		WriteError(w, http.StatusBadRequest, "execution 'workflow_id' is required")
		return
	}

	exec, err := c.service.CreateExecution(r.Context(), req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]any{
		"status":       "success",
		"execution_id": exec.ID,
		"state":        exec.Status,
	})
}

// Update handles (PUT) execution/update
func (c *ExecutionController) Update(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ExecutionID == "" || req.StepNo <= 0 {
		WriteError(w, http.StatusBadRequest, "'execution_id' and valid 'step_no' are required")
		return
	}

	if err := c.service.UpdateExecution(r.Context(), req); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"status":  "success",
		"step_no": req.StepNo,
	})
}

// List handles (GET) execution
func (c *ExecutionController) List(w http.ResponseWriter, r *http.Request) {
	executions, err := c.repo.List()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if executions == nil {
		executions = []models.Execution{}
	}
	WriteJSON(w, http.StatusOK, executions)
}

// GetByID handles (GET) execution/:id
func (c *ExecutionController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	details, err := c.service.GetExecutionDetails(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "execution not found")
		return
	}
	WriteJSON(w, http.StatusOK, details)
}
