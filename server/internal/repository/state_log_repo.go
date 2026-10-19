package repository

import (
	"database/sql"
	"fmt"
	"time"

	"temporal-clone/server/internal/models"
)

type StateLogRepo struct {
	db *sql.DB
}

func NewStateLogRepo(db *sql.DB) *StateLogRepo {
	return &StateLogRepo{db: db}
}

func (r *StateLogRepo) Append(log *models.StateLog) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now().UTC()
	}

	query := `
	INSERT INTO state_logs (execution_id, step_no, activity_id, event, state, timestamp)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, log.ExecutionID, log.StepNo, log.ActivityID, log.Event, log.State, log.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to append state log: %w", err)
	}
	return nil
}

func (r *StateLogRepo) GetByExecutionID(executionID string) ([]models.StateLog, error) {
	query := `
	SELECT execution_id, step_no, activity_id, event, state, timestamp
	FROM state_logs
	WHERE execution_id = $1
	ORDER BY step_no ASC
	`
	rows, err := r.db.Query(query, executionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.StateLog
	for rows.Next() {
		var l models.StateLog
		var actID sql.NullString
		if err := rows.Scan(&l.ExecutionID, &l.StepNo, &actID, &l.Event, &l.State, &l.Timestamp); err != nil {
			return nil, err
		}
		if actID.Valid {
			l.ActivityID = &actID.String
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *StateLogRepo) GetLatestStepNo(executionID string) (int, error) {
	query := `
	SELECT COALESCE(MAX(step_no), 0)
	FROM state_logs
	WHERE execution_id = $1
	`
	var latestStep int
	err := r.db.QueryRow(query, executionID).Scan(&latestStep)
	if err != nil {
		return 0, err
	}
	return latestStep, nil
}
