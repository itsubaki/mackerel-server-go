package usecase

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (s *UserInteractor) List() (domain.Users, error) {
	return s.UserRepository.FindAll()
}

func (s *UserInteractor) Delete(userId string) error {
	return s.UserRepository.Delete(userId)
}
