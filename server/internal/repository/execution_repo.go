package repository

import (
	"database/sql"
	"fmt"
	"time"

	"temporal-clone/server/internal/models"
)

type ExecutionRepo struct {
	db *sql.DB
}

func NewExecutionRepo(db *sql.DB) *ExecutionRepo {
	return &ExecutionRepo{db: db}
}

func (r *ExecutionRepo) Create(exec *models.Execution) error {
	now := time.Now().UTC()
	exec.CreatedAt = now
	exec.UpdatedAt = now
	if exec.Status == "" {
		exec.Status = models.ExecutionStatusRunning
	}

	query := `
	INSERT INTO executions (id, workflow_id, status, created_at, updated_at, completed_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, exec.ID, exec.WorkflowID, exec.Status, exec.CreatedAt, exec.UpdatedAt, exec.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to create execution: %w", err)
	}
	return nil
}

func (r *ExecutionRepo) UpdateStatus(id string, status models.ExecutionStatus, completedAt *time.Time) error {
	now := time.Now().UTC()
	query := `
	UPDATE executions 
	SET status = $1, updated_at = $2, completed_at = $3
	WHERE id = $4
	`
	result, err := r.db.Exec(query, status, now, completedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update execution status: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ExecutionRepo) GetByID(id string) (*models.Execution, error) {
	query := `
	SELECT id, workflow_id, status, created_at, updated_at, completed_at
	FROM executions
	WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var exec models.Execution
	var completedAt sql.NullTime
	err := row.Scan(&exec.ID, &exec.WorkflowID, &exec.Status, &exec.CreatedAt, &exec.UpdatedAt, &completedAt)
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		exec.CompletedAt = &completedAt.Time
	}
	return &exec, nil
}

func (r *ExecutionRepo) List() ([]models.Execution, error) {
	query := `
	SELECT id, workflow_id, status, created_at, updated_at, completed_at
	FROM executions
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var executions []models.Execution
	for rows.Next() {
		var exec models.Execution
		var completedAt sql.NullTime
		if err := rows.Scan(&exec.ID, &exec.WorkflowID, &exec.Status, &exec.CreatedAt, &exec.UpdatedAt, &completedAt); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			exec.CompletedAt = &completedAt.Time
		}
		executions = append(executions, exec)
	}
	return executions, rows.Err()
}
