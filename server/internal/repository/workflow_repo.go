package repository

import (
	"database/sql"
	"fmt"
	"time"

	"temporal-clone/server/internal/models"
)

type WorkflowRepo struct {
	db *sql.DB
}

func NewWorkflowRepo(db *sql.DB) *WorkflowRepo {
	return &WorkflowRepo{db: db}
}

func (r *WorkflowRepo) Create(wf *models.Workflow) error {
	now := time.Now().UTC()
	wf.CreatedAt = now
	wf.UpdatedAt = now

	query := `
	INSERT INTO workflows (id, name, description, version, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, wf.ID, wf.Name, wf.Description, wf.Version, wf.CreatedAt, wf.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}
	return nil
}

func (r *WorkflowRepo) Update(wf *models.Workflow) error {
	wf.UpdatedAt = time.Now().UTC()

	query := `
	UPDATE workflows 
	SET name = $1, description = $2, version = $3, updated_at = $4
	WHERE id = $5
	`
	result, err := r.db.Exec(query, wf.Name, wf.Description, wf.Version, wf.UpdatedAt, wf.ID)
	if err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
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

func (r *WorkflowRepo) GetByID(id string) (*models.Workflow, error) {
	query := `
	SELECT id, name, description, version, created_at, updated_at
	FROM workflows
	WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var wf models.Workflow
	err := row.Scan(&wf.ID, &wf.Name, &wf.Description, &wf.Version, &wf.CreatedAt, &wf.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wf, nil
}

func (r *WorkflowRepo) List() ([]models.Workflow, error) {
	query := `
	SELECT id, name, description, version, created_at, updated_at
	FROM workflows
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.Workflow
	for rows.Next() {
		var wf models.Workflow
		if err := rows.Scan(&wf.ID, &wf.Name, &wf.Description, &wf.Version, &wf.CreatedAt, &wf.UpdatedAt); err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}
	return workflows, rows.Err()
}
