package repository

import "database/sql"

type TenantConfig struct {
	TenantId          string
	MaxJobs           int
	MaxConcurrentJobs int
	RateLimitPerSec   int
	Priority          int
}

type TenantConfigRepo struct {
	db *sql.DB
}

func NewTenantConfigRepo(db *sql.DB) *TenantConfigRepo {
	return &TenantConfigRepo{db: db}
}

func (r *TenantConfigRepo) UpsertConfig(c *TenantConfig) error {
	_, err := r.db.Exec(`
		INSERT INTO tenant_configs
		(tenant_id, max_jobs, max_concurrent_jobs, rate_limit_per_sec, priority)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tenant_id)
		DO UPDATE SET
		max_jobs=$2,
		max_concurrent_jobs=$3,
		rate_limit_per_sec=$4,
		priority=$5
	`, c.TenantId, c.MaxJobs, c.MaxConcurrentJobs, c.RateLimitPerSec, c.Priority)

	return err
}

func (r *TenantConfigRepo) GetConfig(tenantId string) (*TenantConfig, error) {
	row := r.db.QueryRow(`
		SELECT tenant_id, max_jobs, max_concurrent_jobs, rate_limit_per_sec, priority
		FROM tenant_configs WHERE tenant_id=$1
	`, tenantId)

	c := &TenantConfig{}
	err := row.Scan(
		&c.TenantId,
		&c.MaxJobs,
		&c.MaxConcurrentJobs,
		&c.RateLimitPerSec,
		&c.Priority,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *TenantConfigRepo) CountJobs(tenantId string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM jobs WHERE tenant_id=$1 AND status='active'
	`, tenantId).Scan(&count)
	return count, err
}

func (r *TenantConfigRepo) CountRunningJobs(tenantId string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM executions
		WHERE tenant_id=$1 AND status='running'
	`, tenantId).Scan(&count)
	return count, err
}
