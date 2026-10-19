package controller

import (
	"encoding/json"
	"net/http"

	"temporal-clone/server/internal/models"
	"temporal-clone/server/internal/repository"
)

type ActivityController struct {
	repo *repository.ActivityRepo
}

func NewActivityController(repo *repository.ActivityRepo) *ActivityController {
	return &ActivityController{repo: repo}
}

// Create handles (PUT) activity/create
func (c *ActivityController) Create(w http.ResponseWriter, r *http.Request) {
	var act models.Activity
	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if act.ID == "" || act.Name == "" {
		WriteError(w, http.StatusBadRequest, "activity 'id' and 'name' are required")
		return
	}
	if err := c.repo.Create(&act); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]any{"status": "success", "data": act})
}

// Update handles (PUT) activity/update
func (c *ActivityController) Update(w http.ResponseWriter, r *http.Request) {
	var act models.Activity
	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if act.ID == "" {
		WriteError(w, http.StatusBadRequest, "activity 'id' is required")
		return
	}
	if err := c.repo.Update(&act); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"status": "success", "data": act})
}

// List handles (GET) activity
func (c *ActivityController) List(w http.ResponseWriter, r *http.Request) {
	activities, err := c.repo.List()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if activities == nil {
		activities = []models.Activity{}
	}
	WriteJSON(w, http.StatusOK, activities)
}

// GetByID handles (GET) activity/:id
func (c *ActivityController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	act, err := c.repo.GetByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "activity not found")
		return
	}
	WriteJSON(w, http.StatusOK, act)
}
