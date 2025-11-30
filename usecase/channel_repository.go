package usecase

import "github.com/itsubaki/mackerel-server-go/domain"

type ChannelRepository interface {
	List(orgID string) (*domain.Channels, error)
	Save(orgID string, channel *domain.Channel) (any, error)
	Exists(orgID, channelID string) bool
	Delete(orgID, channelID string) (any, error)
}
