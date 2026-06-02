package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	AdminUsername string
	AdminPassword string
	R2AccountID   string
	R2AccessKey   string
	R2SecretKey   string
	R2BucketName  string
	R2PublicURL   string
	AllowedOrigins []string
}

var C *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	C = &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   mustEnv("DATABASE_URL"),
		JWTSecret:     mustEnv("JWT_SECRET"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: mustEnv("ADMIN_PASSWORD"),
		R2AccountID:   getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKey:   getEnv("R2_ACCESS_KEY", ""),
		R2SecretKey:   getEnv("R2_SECRET_KEY", ""),
		R2BucketName:  getEnv("R2_BUCKET_NAME", ""),
		R2PublicURL:   getEnv("R2_PUBLIC_URL", ""),
		AllowedOrigins: []string{
			getEnv("FRONTEND_URL", "http://localhost:4321"),
			getEnv("ADMIN_URL", "http://localhost:5173"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}
