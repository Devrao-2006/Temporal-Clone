package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"temporal-clone/server/internal/controller"
	"temporal-clone/server/internal/repository"
	"temporal-clone/server/internal/routes"
	"temporal-clone/server/internal/service"
)

func main() {
	port := flag.String("port", getEnv("PORT", "8080"), "Port for the HTTP server to listen on")
	dbURL := flag.String("db", getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/temporal?sslmode=disable"), "PostgreSQL connection string")
	flag.Parse()

	log.Printf("[Server] Initializing PostgreSQL database...")
	db, err := repository.InitDB(*dbURL)
	if err != nil {
		log.Fatalf("[Server] Failed to initialize PostgreSQL database: %v", err)
	}
	defer db.Close()
	log.Printf("[Server] PostgreSQL database initialized successfully (all 4 tables ready).")

	// Repositories (Database Layer)
	workflowRepo := repository.NewWorkflowRepo(db)
	activityRepo := repository.NewActivityRepo(db)
	executionRepo := repository.NewExecutionRepo(db)
	stateLogRepo := repository.NewStateLogRepo(db)

	// Services (Business Logic)
	executionService := service.NewExecutionService(
		workflowRepo,
		activityRepo,
		executionRepo,
		stateLogRepo,
	)

	// Controllers (HTTP Handlers)
	wfController := controller.NewWorkflowController(workflowRepo)
	actController := controller.NewActivityController(activityRepo)
	execController := controller.NewExecutionController(executionRepo, executionService)

	// Routes (Routing Layer)
	mux := http.NewServeMux()
	routes.Register(mux, wfController, actController, execController)

	serverAddr := ":" + *port
	log.Printf("[Server] Starting HTTP server on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("[Server] Server failed: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
