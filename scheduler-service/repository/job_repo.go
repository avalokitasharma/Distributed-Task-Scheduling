package repository

import (
	"database/sql"
	"time"
)

type Job struct {
	ID              string
	TenantID        string
	Name            string
	TaskType        string
	Params          []byte
	ScheduleType    string // ONCE | CRON | INTERVAL
	CronExpression  string
	IntervalSeconds int
	NextRunAt       time.Time
	Status          string
	RetryPolicy     []byte
	CreatedAt       time.Time
}

type JobsRepo struct {
	db *sql.DB
}

func NewJobRepo(db *sql.DB) *JobsRepo {
	return &JobsRepo{db: db}
}

func (r *JobsRepo) Create(j *Job) error {
	_, err := r.db.Exec(`
		INSERT INTO JOBS
		(id, tenant_id, name, task_type, params, schedule_type, cron_expression, internval_seconds, nex_run_at, status, retry_policy, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, j.ID,
		j.TenantID,
		j.Name,
		j.TaskType,
		j.Params,
		j.ScheduleType,
		j.CronExpression,
		j.IntervalSeconds,
		j.NextRunAt,
		j.Status,
		j.RetryPolicy,
		time.Now(),
	)
	return err
}

func (r *JobsRepo) UpdateNextRun(jobID string, next time.Time) error {
	_, err := r.db.Exec(`
		UPDATE jobs SET next_run_at=$1 WHERE id=$2
	`, next, jobID)

	return err
}
