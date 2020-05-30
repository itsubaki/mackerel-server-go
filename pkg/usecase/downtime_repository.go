package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type DowntimeRepository interface {
	List(orgID string) (*domain.Downtimes, error)
	Save(orgID string, downtime *domain.Downtime) (*domain.Downtime, error)
	Update(orgID string, downtime *domain.Downtime) (*domain.Downtime, error)
	Downtime(orgID, downtimeID string) (*domain.Downtime, error)
	Delete(orgID, downtimeID string) (*domain.Downtime, error)
}
