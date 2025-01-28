package config

import "os"

// GetEnv retrieves environment variables with a fallback default
func GetEnv(key string, fallback string) string {
    val := os.Getenv(key)
    if val == "" {
        return fallback
    }
    return val
}