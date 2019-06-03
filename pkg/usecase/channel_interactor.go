package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ChannelInteractor struct {
	ChannelRepository ChannelRepository
}

func (s *ChannelInteractor) List(orgID string) (*domain.Channels, error) {
	return s.ChannelRepository.List(orgID)
}

func (s *ChannelInteractor) Save(orgID string, channel *domain.Channel) (*domain.Channel, error) {
	return s.ChannelRepository.Save(orgID, channel)
}

func (s *ChannelInteractor) Exists(orgID, channelID string) bool {
	return s.ChannelRepository.Exists(orgID, channelID)
}

func (s *ChannelInteractor) Delete(orgID, channelID string) (*domain.Channel, error) {
	return s.ChannelRepository.Delete(orgID, channelID)
}
