package config

import "os"

type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	JWTAccessSecret  string
	JWTRefreshSecret string
}

func Load() Config {
	return Config{
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/app?sslmode=disable"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
