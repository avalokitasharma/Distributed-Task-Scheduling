package service

import (
	"github.com/avalokitasharma/job-scheduler/scheduler-service/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SchedulerService struct {
	jobRepo    *repository.JobsRepo
	redis      *redis.Client
	tenantHost string
}

func NewSchedulerService(r *repository.JobsRepo, rc *redis.Client, tenantHost string) *SchedulerService {
	return &SchedulerService{
		jobRepo:    r,
		redis:      rc,
		tenantHost: tenantHost,
	}
}

// create job

func (s *SchedulerService) CreateJob(j *repository.Job) error {
	// 1. Quota check
	if err := s.checkQuota(j.TenantID); err != nil {
		return err
	}
	// 2. Compute next run
	nextRun, err := computeNextRun(j)
	if err != nil {
		return err
	}
	j.NextRunAt = nextRun
	j.ID = uuid.NewString()
	j.Status = "active"

	// 3. Add to DB
	if err := s.jobRepo.Create(j); err != nil {
		return err
	}

	// 4. Push to Redis sorted set
	return nil
}

func (s *SchedulerService) checkQuota(tenantId string) error {

	return nil
}
