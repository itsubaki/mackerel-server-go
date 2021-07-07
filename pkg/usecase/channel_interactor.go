package usecase

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type ChannelInteractor struct {
	ChannelRepository ChannelRepository
}

func (intr *ChannelInteractor) List(orgID string) (*domain.Channels, error) {
	return intr.ChannelRepository.List(orgID)
}

func (intr *ChannelInteractor) Save(orgID string, channel *domain.Channel) (interface{}, error) {
	channel.ID = domain.NewRandomID()
	return intr.ChannelRepository.Save(orgID, channel)
}

func (intr *ChannelInteractor) Exists(orgID, channelID string) bool {
	return intr.ChannelRepository.Exists(orgID, channelID)
}

func (intr *ChannelInteractor) Delete(orgID, channelID string) (interface{}, error) {
	if !intr.ChannelRepository.Exists(orgID, channelID) {
		return nil, &ChannelNotFound{Err{errors.New(fmt.Sprintf("when the supported channel can not be found in <%s>", channelID))}}
	}

	return intr.ChannelRepository.Delete(orgID, channelID)
}
