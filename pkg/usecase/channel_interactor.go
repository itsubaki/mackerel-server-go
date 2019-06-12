package usecase

import (
	"errors"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ChannelInteractor struct {
	ChannelRepository ChannelRepository
}

func (s *ChannelInteractor) List(orgID string) (*domain.Channels, error) {
	return s.ChannelRepository.List(orgID)
}

func (s *ChannelInteractor) Save(orgID string, channel *domain.Channel) (interface{}, error) {
	channel.ID = domain.NewRandomID(11)
	return s.ChannelRepository.Save(orgID, channel)
}

func (s *ChannelInteractor) Exists(orgID, channelID string) bool {
	return s.ChannelRepository.Exists(orgID, channelID)
}

func (s *ChannelInteractor) Delete(orgID, channelID string) (interface{}, error) {
	if !s.ChannelRepository.Exists(orgID, channelID) {
		return nil, &ChannelNotFound{Err{errors.New("when the supported channel can not be found in <channelId>")}}
	}

	return s.ChannelRepository.Delete(orgID, channelID)
}
