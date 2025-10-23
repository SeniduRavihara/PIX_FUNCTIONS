package utils

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	Environment    string
	FirecrackerBin string
	KernelPath     string
	RootFSPath     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "host=localhost user=voltrun password=voltrun dbname=voltrun port=5432 sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "voltrun-secret-change-in-production"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		FirecrackerBin: getEnv("FIRECRACKER_BIN", "/usr/bin/firecracker"),
		KernelPath:     getEnv("KERNEL_PATH", "/var/lib/voltrun/vmlinux.bin"),
		RootFSPath:     getEnv("ROOTFS_PATH", "/var/lib/voltrun/rootfs.ext4"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool gets an environment variable as bool or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
