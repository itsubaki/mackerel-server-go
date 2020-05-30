package usecase

import "github.com/itsubaki/mackerel-server-go/pkg/domain"

type ChannelRepository interface {
	List(orgID string) (*domain.Channels, error)
	Save(orgID string, channel *domain.Channel) (interface{}, error)
	Exists(orgID, channelID string) bool
	Delete(orgID, channelID string) (interface{}, error)
}
