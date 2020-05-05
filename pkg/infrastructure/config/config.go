package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	Driver   string
	Host     string
	Database string
}

func GetValue(envKey, defaultValue string) string {
	if len(os.Getenv(envKey)) > 0 {
		return os.Getenv(envKey)
	}

	return defaultValue
}

func New() *Config {
	return &Config{
		Port:     GetValue("PORT", ":8080"),
		Driver:   GetValue("DRIVER", "mysql"),
		Host:     GetValue("HOST", "root:secret@tcp(127.0.0.1:3306)/"),
		Database: GetValue("DATABASE", "mackerel"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s%s", c.Host, c.Database)
}
