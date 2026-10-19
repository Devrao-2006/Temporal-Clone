package routes

import (
	"net/http"

	"temporal-clone/server/internal/controller"
)

func Register(
	mux *http.ServeMux,
	wfCtrl *controller.WorkflowController,
	actCtrl *controller.ActivityController,
	execCtrl *controller.ExecutionController,
) {
	// Workflow APIs
	mux.HandleFunc("PUT /workflow/create", wfCtrl.Create)
	mux.HandleFunc("PUT /workflow/update", wfCtrl.Update)
	mux.HandleFunc("GET /workflow", wfCtrl.List)
	mux.HandleFunc("GET /workflow/{id}", wfCtrl.GetByID)

	// Activity APIs
	mux.HandleFunc("PUT /activity/create", actCtrl.Create)
	mux.HandleFunc("PUT /activity/update", actCtrl.Update)
	mux.HandleFunc("GET /activity", actCtrl.List)
	mux.HandleFunc("GET /activity/{id}", actCtrl.GetByID)

	// Execution APIs
	mux.HandleFunc("PUT /execution/create", execCtrl.Create)
	mux.HandleFunc("PUT /execution/update", execCtrl.Update)
	mux.HandleFunc("GET /execution", execCtrl.List)
	mux.HandleFunc("GET /execution/{id}", execCtrl.GetByID)
}
