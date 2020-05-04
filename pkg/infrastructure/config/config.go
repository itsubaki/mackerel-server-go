package config

import "os"

type Config struct {
	Port           string
	Driver         string
	DataSourceName string
	DatabaseName   string
}

func GetValue(envKey, defaultValue string) string {
	if len(os.Getenv(envKey)) > 0 {
		return os.Getenv(envKey)
	}

	return defaultValue
}

func New() *Config {
	return &Config{
		Port:           GetValue("PORT", ":8080"),
		Driver:         GetValue("DRIVER", "mysql"),
		DataSourceName: GetValue("DATA_SOURCE_NAME", "root:secret@tcp(127.0.0.1:3306)/"),
		DatabaseName:   GetValue("DATABASE_NAME", "mackerel"),
	}
}
