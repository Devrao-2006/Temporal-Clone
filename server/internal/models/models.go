package models

import "time"

// Workflow represents a workflow definition (Table 1)
type Workflow struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Activity represents an activity definition (Table 2)
type Activity struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	InputSchema  string    `json:"input_schema"`
	OutputSchema string    `json:"output_schema"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ExecutionStatus represents the lifecycle state of an execution
type ExecutionStatus string

const (
	ExecutionStatusRunning   ExecutionStatus = "RUNNING"
	ExecutionStatusCompleted ExecutionStatus = "COMPLETED"
	ExecutionStatusFailed    ExecutionStatus = "FAILED"
)

// Execution represents a running or finished workflow execution (Table 3)
type Execution struct {
	ID          string          `json:"id"`
	WorkflowID  string          `json:"workflow_id"`
	Status      ExecutionStatus `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

// StateLog represents an individual step event history record (Table 4)
type StateLog struct {
	ExecutionID string    `json:"execution_id"`
	StepNo      int       `json:"step_no"`
	ActivityID  *string   `json:"activity_id,omitempty"`
	Event       string    `json:"event"`
	State       string    `json:"state"` // JSON string or text representing current state/output
	Timestamp   time.Time `json:"timestamp"`
}

// ExecutionDetails contains execution status along with full history logs
type ExecutionDetails struct {
	Execution Execution  `json:"execution"`
	History   []StateLog `json:"history"`
}
