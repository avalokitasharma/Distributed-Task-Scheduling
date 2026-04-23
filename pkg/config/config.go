package config

import "os"

type Config struct {
	DBURL     string
	RedisAddr string
}

var App *Config

func Load() {
	App = &Config{
		DBURL:     get("DB_URL", "postgres://postgres:postgres@db:5432/jobs?sslmode=disable"),
		RedisAddr: get("REDIS_ADDR", "redis:6379"),
	}
}

func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
