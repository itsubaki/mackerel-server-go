package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ChannelRepository interface {
	List(orgID string) (*domain.Channels, error)
	Save(orgID string, channel *domain.Channel) (*domain.Channel, error)
	Exists(orgID, channelID string) bool
	Delete(orgID, channelID string) (*domain.Channel, error)
}
