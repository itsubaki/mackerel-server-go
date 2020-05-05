package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port     string
	Driver   string
	Host     string
	Database string
}

func GetValue(key, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) > 0 {
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
	}
}

func (c *Config) DSN() string {
	if !strings.HasSuffix(c.Host, "/") && !strings.HasPrefix(c.Database, "/") {
		return fmt.Sprintf("%s/%s", c.Host, c.Database)
	}

	return fmt.Sprintf("%s%s", c.Host, c.Database)
}
