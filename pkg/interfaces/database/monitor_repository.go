package database

type MonitorRepository struct {
	SQLHandler
}

func NewMonitorRepository(handler SQLHandler) *MonitorRepository {
	return &MonitorRepository{
		SQLHandler: handler,
	}
}
