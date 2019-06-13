package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type NotificationGroupInteractor struct {
	NotificationGroupRepository NotificationGroupRepository
}

func (s *NotificationGroupInteractor) List(orgID string) (*domain.NotificationGroups, error) {
	return s.NotificationGroupRepository.List(orgID)
}

func (s *NotificationGroupInteractor) Save(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	group.ID = domain.NewRandomID(11)
	return s.NotificationGroupRepository.Save(orgID, group)
}

func (s *NotificationGroupInteractor) Update(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	if !s.NotificationGroupRepository.Exists(orgID, group.ID) {
		return nil, &NotificationGroupNotFound{Err{errors.New("when the specified notification group does not exist")}}
	}

	return s.NotificationGroupRepository.Update(orgID, group)
}

func (s *NotificationGroupInteractor) Delete(orgID, groupID string) (*domain.NotificationGroup, error) {
	if !s.NotificationGroupRepository.Exists(orgID, groupID) {
		return nil, &NotificationGroupNotFound{Err{errors.New("when the specified notification group does not exist")}}
	}

	return s.NotificationGroupRepository.Delete(orgID, groupID)
}
