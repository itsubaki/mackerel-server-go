package config

import (
	"os"
	"time"
)

type Config struct {
	Port     string
	Driver   string
	Host     string
	Database string
	SQLMode  string
	Timeout  time.Duration
	Sleep    time.Duration
}

func GetValue(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defaultValue
}

func New() *Config {
	return &Config{
		Port:     GetValue("PORT", ":8080"),
		Driver:   GetValue("DRIVER", "mysql"),
		Host:     GetValue("HOST", "root:secret@tcp(127.0.0.1:3306)/"),
		Database: GetValue("DATABASE", "mackerel"),
		SQLMode:  GetValue("SQL_MODE", "release"),
		Timeout:  10 * time.Minute,
		Sleep:    10 * time.Second,
	}
}
