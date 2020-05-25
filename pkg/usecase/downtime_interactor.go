package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type DowntimeInteractor struct {
	DowntimeRepository DowntimeRepository
}

func (s *DowntimeInteractor) List(orgID string) (*domain.Downtimes, error) {
	return nil, nil
}

func (s *DowntimeInteractor) Save(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	return nil, nil
}

func (s *DowntimeInteractor) Update(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	return nil, nil
}

func (s *DowntimeInteractor) Downtime(orgID, downtimeID string) (*domain.Downtime, error) {
	return nil, nil
}

func (s *DowntimeInteractor) Delete(orgID, downtimeID string) (*domain.Downtime, error) {
	return nil, nil
}
