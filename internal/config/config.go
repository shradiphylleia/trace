package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port string
	BaseURL string
	DatabaseURL string
	RedisAddr string
	MinIOEndpoint string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket string
	MinIOUseSSL bool
	CleanupInterval time.Duration
}

func Load() Config {
	return Config{
		Port:env("PORT", "8080"),
		BaseURL:env("BASE_URL", "http://localhost:8080"),
		DatabaseURL:env("DATABASE_URL","postgres://traceshare:traceshare@localhost:5432/traceshare?sslmode=disable"),
		RedisAddr:env("REDIS_ADDR","localhost:6379"),
		MinIOEndpoint:env("MINIO_ENDPOINT","localhost:9000"),
		MinIOAccessKey:env("MINIO_ACCESS_KEY","minioadmin"),
		MinIOSecretKey:env("MINIO_SECRET_KEY","minioadmin"),
		MinIOBucket:env("MINIO_BUCKET","traceshare-artifacts"),
		MinIOUseSSL:envBool("MINIO_USE_SSL",false),
		CleanupInterval:envDuration("CLEANUP_INTERVAL",15*time.Minute),
	}
}
func env(key,fallback string) string {
	if value:=os.Getenv(key);value!=""{
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value:=os.Getenv(key)
	if value==""{
		return fallback
	}
	parsed,err:=strconv.ParseBool(value)
	if err!=nil {
		return fallback
	}
	return parsed
}

func envDuration(key string,fallback time.Duration) time.Duration {
	value:=os.Getenv(key)
	if value=="" {
		return fallback
	}
	parsed,err:=time.ParseDuration(value)
	if err!=nil {
		return fallback
	}
	return parsed
}
