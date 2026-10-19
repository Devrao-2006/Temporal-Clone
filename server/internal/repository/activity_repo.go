package repository

import (
	"database/sql"
	"fmt"
	"time"

	"temporal-clone/server/internal/models"
)

type ActivityRepo struct {
	db *sql.DB
}

func NewActivityRepo(db *sql.DB) *ActivityRepo {
	return &ActivityRepo{db: db}
}

func (r *ActivityRepo) Create(act *models.Activity) error {
	now := time.Now().UTC()
	act.CreatedAt = now
	act.UpdatedAt = now

	query := `
	INSERT INTO activities (id, name, description, input_schema, output_schema, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, act.ID, act.Name, act.Description, act.InputSchema, act.OutputSchema, act.CreatedAt, act.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}
	return nil
}

func (r *ActivityRepo) Update(act *models.Activity) error {
	act.UpdatedAt = time.Now().UTC()

	query := `
	UPDATE activities 
	SET name = $1, description = $2, input_schema = $3, output_schema = $4, updated_at = $5
	WHERE id = $6
	`
	result, err := r.db.Exec(query, act.Name, act.Description, act.InputSchema, act.OutputSchema, act.UpdatedAt, act.ID)
	if err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
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

func (r *ActivityRepo) GetByID(id string) (*models.Activity, error) {
	query := `
	SELECT id, name, description, input_schema, output_schema, created_at, updated_at
	FROM activities
	WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var act models.Activity
	err := row.Scan(&act.ID, &act.Name, &act.Description, &act.InputSchema, &act.OutputSchema, &act.CreatedAt, &act.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &act, nil
}

func (r *ActivityRepo) List() ([]models.Activity, error) {
	query := `
	SELECT id, name, description, input_schema, output_schema, created_at, updated_at
	FROM activities
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var act models.Activity
		if err := rows.Scan(&act.ID, &act.Name, &act.Description, &act.InputSchema, &act.OutputSchema, &act.CreatedAt, &act.UpdatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, act)
	}
	return activities, rows.Err()
}
