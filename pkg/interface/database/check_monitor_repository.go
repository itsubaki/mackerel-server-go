package database

type CheckMonitorRepository struct {
	SQLHandler
}

func NewCheckMonitorRepository(handler SQLHandler) *CheckMonitorRepository {
	return &CheckMonitorRepository{
		SQLHandler: handler,
	}
}
