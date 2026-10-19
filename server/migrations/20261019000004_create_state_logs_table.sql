-- Table 4: State Log Table
CREATE TABLE IF NOT EXISTS state_logs (
    execution_id TEXT NOT NULL REFERENCES executions(id),
    step_no INTEGER NOT NULL,
    activity_id TEXT,
    event TEXT NOT NULL,
    state TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (execution_id, step_no)
);

-- Index for faster retrieval of execution history
CREATE INDEX IF NOT EXISTS idx_state_logs_exec_step ON state_logs (execution_id, step_no);
