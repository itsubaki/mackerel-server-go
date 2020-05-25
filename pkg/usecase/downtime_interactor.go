package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type DowntimeInteractor struct {
	DowntimeRepository DowntimeRepository
}

func (s *DowntimeInteractor) List(orgID string) (*domain.Downtimes, error) {
	return s.DowntimeRepository.List(orgID)
}

func (s *DowntimeInteractor) Save(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	downtime.ID = domain.NewRandomID(11)
	return s.DowntimeRepository.Save(orgID, downtime)
}

func (s *DowntimeInteractor) Update(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	return s.DowntimeRepository.Update(orgID, downtime)
}

func (s *DowntimeInteractor) Downtime(orgID, downtimeID string) (*domain.Downtime, error) {
	return s.DowntimeRepository.Downtime(orgID, downtimeID)
}

func (s *DowntimeInteractor) Delete(orgID, downtimeID string) (*domain.Downtime, error) {
	return s.DowntimeRepository.Delete(orgID, downtimeID)
}
