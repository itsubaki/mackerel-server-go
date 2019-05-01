package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type UserInteractor struct {
	UserRepository UserRepository
}

func (s *UserInteractor) List() (*domain.Users, error) {
	return s.UserRepository.List()
}

func (s *UserInteractor) Delete(userID string) (*domain.User, error) {
	return s.UserRepository.Delete(userID)
}
