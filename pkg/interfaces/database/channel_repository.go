package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ChannelRepository struct {
	SQLHandler
}

func NewChannelRepository(handler SQLHandler) *ChannelRepository {
	return &ChannelRepository{
		SQLHandler: handler,
	}
}

func (repo *ChannelRepository) List(orgID string) (*domain.Channels, error) {
	return nil, nil
}

func (repo *ChannelRepository) Save(orgID string, channel *domain.Channel) (*domain.Channel, error) {
	return nil, nil
}

func (repo *ChannelRepository) Exists(orgID, channelID string) bool {
	return true
}

func (repo *ChannelRepository) Delete(orgID, channelID string) (*domain.Channel, error) {
	return nil, nil
}
