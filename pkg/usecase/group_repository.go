package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type NotificationGroupRepository interface {
	List(orgID string) (*domain.NotificationGroups, error)
	Save(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error)
	Exists(orgID, groupID string) bool
	Update(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error)
	Delete(orgID, groupID string) (*domain.NotificationGroup, error)
}
