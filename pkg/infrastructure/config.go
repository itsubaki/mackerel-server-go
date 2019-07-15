package infrastructure

type Config struct {
	Port           string
	Driver         string
	DataSourceName string
	DatabaseName   string
}

func NewConfig() *Config {
	return &Config{
		Port:           ":8080",
		Driver:         "mysql",
		DataSourceName: "root:secret@tcp(127.0.0.1:3307)/",
		DatabaseName:   "mackerel",
	}
}
