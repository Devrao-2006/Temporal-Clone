package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"temporal-clone/server/internal/models"
	"temporal-clone/server/internal/repository"
)

type ExecutionService struct {
	workflowRepo  *repository.WorkflowRepo
	activityRepo  *repository.ActivityRepo
	executionRepo *repository.ExecutionRepo
	stateLogRepo  *repository.StateLogRepo
}

func NewExecutionService(
	wfRepo *repository.WorkflowRepo,
	actRepo *repository.ActivityRepo,
	execRepo *repository.ExecutionRepo,
	logRepo *repository.StateLogRepo,
) *ExecutionService {
	return &ExecutionService{
		workflowRepo:  wfRepo,
		activityRepo:  actRepo,
		executionRepo: execRepo,
		stateLogRepo:  logRepo,
	}
}

type CreateExecutionRequest struct {
	ID         string                 `json:"id"`
	WorkflowID string                 `json:"workflow_id"`
	Input      map[string]interface{} `json:"input"`
}

func (s *ExecutionService) CreateExecution(ctx context.Context, req CreateExecutionRequest) (*models.Execution, error) {
	// 1. Verify workflow exists
	_, err := s.workflowRepo.GetByID(req.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("workflow not found: %w", err)
	}

	execID := req.ID
	if execID == "" {
		execID = uuid.New().String()
	}

	exec := &models.Execution{
		ID:         execID,
		WorkflowID: req.WorkflowID,
		Status:     models.ExecutionStatusRunning,
	}

	// 2. Persist execution in DB (Table 3)
	if err := s.executionRepo.Create(exec); err != nil {
		return nil, fmt.Errorf("failed to create execution record: %w", err)
	}

	// 3. Persist initial state log in DB (Table 4)
	inputBytes, _ := json.Marshal(req.Input)
	initialLog := &models.StateLog{
		ExecutionID: execID,
		StepNo:      1,
		ActivityID:  nil,
		Event:       "INIT-WORKFLOW",
		State:       string(inputBytes),
		Timestamp:   time.Now().UTC(),
	}
	if err := s.stateLogRepo.Append(initialLog); err != nil {
		return nil, fmt.Errorf("failed to record initial state log: %w", err)
	}

	//Send the data to Matching server after this

	return exec, nil
}

type UpdateExecutionRequest struct {
	ExecutionID string                 `json:"execution_id"`
	StepNo      int                    `json:"step_no"`
	ActivityID  *string                `json:"activity_id,omitempty"`
	Event       string                 `json:"event"`
	State       map[string]interface{} `json:"state"`
	Status      models.ExecutionStatus `json:"status"` // RUNNING, COMPLETED, FAILED
}

func (s *ExecutionService) UpdateExecution(ctx context.Context, req UpdateExecutionRequest) error {
	stateBytes, _ := json.Marshal(req.State)
	log := &models.StateLog{
		ExecutionID: req.ExecutionID,
		StepNo:      req.StepNo,
		ActivityID:  req.ActivityID,
		Event:       req.Event,
		State:       string(stateBytes),
		Timestamp:   time.Now().UTC(),
	}

	// 1. Append step to Table 4
	if err := s.stateLogRepo.Append(log); err != nil {
		return fmt.Errorf("failed to append state log: %w", err)
	}

	// 2. Update Table 3 if status changed
	if req.Status == models.ExecutionStatusCompleted || req.Status == models.ExecutionStatusFailed {
		now := time.Now().UTC()
		if err := s.executionRepo.UpdateStatus(req.ExecutionID, req.Status, &now); err != nil {
			return fmt.Errorf("failed to update execution status: %w", err)
		}
	}

	return nil
}

func (s *ExecutionService) GetExecutionDetails(ctx context.Context, id string) (*models.ExecutionDetails, error) {
	exec, err := s.executionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("execution not found: %w", err)
	}

	history, err := s.stateLogRepo.GetByExecutionID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve history: %w", err)
	}

	return &models.ExecutionDetails{
		Execution: *exec,
		History:   history,
	}, nil
}
