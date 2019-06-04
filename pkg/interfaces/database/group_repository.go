package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type NotificationGroupRepository struct {
	SQLHandler
}

func NewNotificationGroupRepository(handler SQLHandler) *NotificationGroupRepository {
	return &NotificationGroupRepository{
		SQLHandler: handler,
	}
}

func (repo *NotificationGroupRepository) List(orgID string) (*domain.NotificationGroups, error) {
	return nil, nil
}

func (repo *NotificationGroupRepository) Save(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	return nil, nil
}

func (repo *NotificationGroupRepository) Exists(orgID, groupID string) bool {
	return true
}

func (repo *NotificationGroupRepository) Update(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	return nil, nil
}

func (repo *NotificationGroupRepository) Delete(orgID, groupID string) (*domain.NotificationGroup, error) {
	return nil, nil
}
