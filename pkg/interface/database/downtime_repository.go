package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type DowntimeRepository struct {
	SQLHandler
}

func NewDowntimeRepository(handler SQLHandler) *DowntimeRepository {
	return &DowntimeRepository{
		SQLHandler: handler,
	}
}

func (repo *DowntimeRepository) List(orgID string) (*domain.Downtimes, error) {
	return nil, nil
}

func (repo *DowntimeRepository) Save(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	return nil, nil
}

func (repo *DowntimeRepository) Update(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	return nil, nil
}

func (repo *DowntimeRepository) Downtime(orgID, downtimeID string) (*domain.Downtime, error) {
	return nil, nil
}

func (repo *DowntimeRepository) Delete(orgID, downtimeID string) (*domain.Downtime, error) {
	return nil, nil
}
