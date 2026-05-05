package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/avalokitasharma/job-scheduler/tenant-service/repository"
	"github.com/redis/go-redis/v9"
)

type TenantConfigService struct {
	repo  *repository.TenantConfigRepo
	redis *redis.Client
}

// ------- Tenant data cached in Redis ---------//
// tenant:{id}:config                 -> JSON (MaxJobs, MaxConcurrentJobs, RateLimit)
// tenant:{id}:jobs:count            -> total jobs
// tenant:{id}:jobs:running          -> running jobs
// tenant:{id}:rate                  -> sorted set (timestamps)

func NewTenantService(r *repository.TenantConfigRepo, rc *redis.Client) *TenantConfigService {
	return &TenantConfigService{
		repo:  r,
		redis: rc,
	}
}

// Get config
func (s *TenantConfigService) GetConfig(ctx context.Context, tenantId string) (*repository.TenantConfig, error) {
	key := "tenant:" + tenantId + ":config"

	//check cache
	val, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var cfg repository.TenantConfig
		json.Unmarshal([]byte(val), &cfg)
		return &cfg, nil
	}
	// Fallback to DB
	cfg, err := s.repo.GetConfig(tenantId)
	if err != nil {
		return nil, err
	}
	// Populate cache
	bytes, _ := json.Marshal(cfg)
	s.redis.Set(ctx, key, bytes, 60*time.Minute)

	return cfg, nil

}

// upsert config
func (s *TenantConfigService) UpsertConfig(ctx context.Context, tenantId string, cfg *repository.TenantConfig) error {
	err := s.repo.UpsertConfig(cfg)
	if err != nil {
		return err
	}
	key := "tenant:" + tenantId + ":config"
	bytes, _ := json.Marshal(cfg)

	// write-through cache
	s.redis.Set(ctx, key, bytes, 60*time.Minute)
	return nil
}

// job creation quota checks
func (s *TenantConfigService) CanCreateJob(ctx context.Context, tenantId string) (bool, error) {
	cfg, err := s.GetConfig(ctx, tenantId)

	key := "tenant:" + tenantId + ":jobs:count"

	count, err := s.redis.Get(ctx, key).Int()
	// fallback to DB
	if err != nil {
		count, err = s.repo.CountJobs(tenantId)
		if err != nil {
			return false, err
		}
	}
	if count >= cfg.MaxJobs {
		return false, errors.New("job quota exceeded")
	}
	return true, nil
}

// Concurrent execution quota
func (s *TenantConfigService) CanRunJob(tenantId string) (bool, error) {
	return true, nil
}

// Rate limiting (Token bucket using Redis)
func (s *TenantConfigService) CheckRateLimit(tenantId string) error {
	return nil
}
