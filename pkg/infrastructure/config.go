package infrastructure

import "fmt"

type Config struct {
	Port           string
	Driver         string
	DataSourceName string
	DatabaseName   string
}

func (c *Config) DataSourceNameWithDatabaseName() string {
	return fmt.Sprintf("%s%s", c.DataSourceName, c.DatabaseName)
}

func (c *Config) QueryCreateDatabase() string {
	return fmt.Sprintf("create database if not exists %s", c.DatabaseName)
}

func NewConfig() *Config {
	return &Config{
		Port:           ":8080",
		Driver:         "mysql",
		DataSourceName: "root:secret@tcp(127.0.0.1:3307)/",
		DatabaseName:   "mackerel",
	}
}
