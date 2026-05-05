package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
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
func (s *TenantConfigService) CanCreateJob(ctx context.Context, tenantId string) error {
	cfg, err := s.GetConfig(ctx, tenantId)

	key := "tenant:" + tenantId + ":jobs:count"

	count, err := s.redis.Get(ctx, key).Int()
	// fallback to DB
	if err != nil {
		count, err = s.repo.CountJobs(tenantId)
		if err != nil {
			return err
		}
	}
	if count >= cfg.MaxJobs {
		return errors.New("job quota exceeded")
	}
	return nil
}

// Concurrent execution quota - critical path
func (s *TenantConfigService) CanRunJob(ctx context.Context, tenantId string) error {
	cfg, err := s.GetConfig(ctx, tenantId)
	if err != nil {
		return err
	}
	key := "tenant:" + tenantId + ":jobs:running"

	running, err := s.redis.Get(ctx, key).Int()
	// fallback to DB
	if err != nil {

		running, err = s.repo.CountRunningJobs(tenantId)
		if err != nil {
			return err
		}
	}
	if running >= cfg.MaxConcurrentJobs {
		return errors.New("concurrent execution limit reached")
	}
	return nil
}

// Rate limiting (Token bucket using Redis)
func (s *TenantConfigService) CheckRateLimit(ctx context.Context, tenantId string) error {
	cfg, err := s.GetConfig(ctx, tenantId)
	if err != nil {
		return err
	}

	key := "tenant:" + tenantId + ":rate"
	now := time.Now().Unix()

	// remove old entries
	s.redis.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now-1, 10))

	count, err := s.redis.ZCard(ctx, key).Result()
	if err != nil {
		return err
	}

	if int(count) >= cfg.RateLimitPerSec {
		return errors.New("rate limit exceeded")
	}

	s.redis.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: now,
	})

	s.redis.Expire(ctx, key, 2*time.Second)
	return nil
	// todo: atomic check - current check then incr causes race condition
}
