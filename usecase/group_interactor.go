package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-server-go/domain"
)

var ErrNotificationGroupNotFound = errors.New("notification group not found")

type NotificationGroupInteractor struct {
	NotificationGroupRepository NotificationGroupRepository
}

func (intr *NotificationGroupInteractor) List(orgID string) (*domain.NotificationGroups, error) {
	return intr.NotificationGroupRepository.List(orgID)
}

func (intr *NotificationGroupInteractor) Save(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	group.ID = domain.NewRandomID()
	return intr.NotificationGroupRepository.Save(orgID, group)
}

func (intr *NotificationGroupInteractor) Update(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	if !intr.NotificationGroupRepository.Exists(orgID, group.ID) {
		return nil, &NotificationGroupNotFound{
			Err{
				Err: ErrNotificationGroupNotFound,
			},
		}
	}

	return intr.NotificationGroupRepository.Update(orgID, group)
}

func (intr *NotificationGroupInteractor) Delete(orgID, groupID string) (*domain.NotificationGroup, error) {
	if !intr.NotificationGroupRepository.Exists(orgID, groupID) {
		return nil, &NotificationGroupNotFound{
			Err{
				Err: ErrNotificationGroupNotFound,
			},
		}
	}

	return intr.NotificationGroupRepository.Delete(orgID, groupID)
}
