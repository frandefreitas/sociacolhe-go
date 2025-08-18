package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	AllowSignup bool
}

func Load() Config {
	c := Config{
		Port:        get("PORT", "8080"),
		DatabaseURL: get("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/sociacolhe?sslmode=disable"),
		JWTSecret:   get("JWT_SECRET", "dev-secret"),
		AllowSignup: get("ALLOW_SIGNUP", "true") == "true",
	}
	log.Printf("config loaded: port=%s allowSignup=%v", c.Port, c.AllowSignup)
	return c
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
