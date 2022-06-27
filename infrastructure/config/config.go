package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	Port            string
	Driver          string
	Host            string
	Database        string
	SQLMode         string
	Timeout         time.Duration
	Sleep           time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	RunFixture      bool
}

func GetValue(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}

func New() *Config {
	runFixture := false
	if strings.ToLower(GetValue("RUN_FIXTURE", "false")) == "true" {
		runFixture = true
	}

	return &Config{
		Port:            GetValue("PORT", "8080"),
		Driver:          GetValue("DRIVER", "mysql"),
		Host:            GetValue("HOST", "root:secret@tcp(127.0.0.1:3306)/"),
		Database:        GetValue("DATABASE", "mackerel"),
		SQLMode:         GetValue("SQL_MODE", "release"),
		Timeout:         10 * time.Minute,
		Sleep:           10 * time.Second,
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxLifetime: 10 * time.Minute,
		RunFixture:      runFixture,
	}
}
