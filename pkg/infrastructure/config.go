package infrastructure

import "os"

type Config struct {
	Port           string
	Driver         string
	DataSourceName string
	DatabaseName   string
}

func GetValue(defaultValue, envKey string) string {
	if len(os.Getenv(envKey)) > 0 {
		return os.Getenv(envKey)
	}

	return defaultValue
}

func NewConfig() *Config {
	return &Config{
		Port:           GetValue(":8080", "PORT"),
		Driver:         GetValue("mysql", "DRIVER"),
		DataSourceName: GetValue("root:secret@tcp(127.0.0.1:3306)/", "DATA_SOURCE_NAME"),
		DatabaseName:   GetValue("mackerel", "DATABASE_NAME"),
	}
}
