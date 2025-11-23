package config

import (
	"os"
	"time"
)

func getTestEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func LoadTestDBConfig() *DBConfig {
	return &DBConfig{
		DBUser:        getTestEnv("TEST_DB_USER", "test"),
		DBPassword:    getTestEnv("TEST_DB_PASSWORD", "test"),
		DBName:        getTestEnv("TEST_DB_NAME", "testdb"),
		DBHost:        getTestEnv("TEST_DB_HOST", "localhost"),
		DBPort:        getTestEnv("TEST_DB_PORT", "5439"),
		DBMaxCons:     5,
		DBMinCons:     1,
		DBConLifetime: time.Hour,
		DBConnTimeout: 5 * time.Second,
	}
}
