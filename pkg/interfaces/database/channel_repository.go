package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ChannelRepository struct {
	SQLHandler
}

func NewChannelRepository(handler SQLHandler) *ChannelRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists channels (
				org_id              varchar(64) not null,
				id                  varchar(16) not null primary key,
				name                varchar(16) not null,
				type                varchar(16) not null,
				url                 text,
				mentions            text,
				enabled_graph_image bool not null default '1'
			)
			`,
		); err != nil {
			return fmt.Errorf("create tables channels: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists channel_events (
				org_id     varchar(64) not null,
				channel_id varchar(16) not null,
				event      varchar(16) not null,
				primary key(channel_id, event)
			)
			`,
		); err != nil {
			return fmt.Errorf("create tables channel_events: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists channel_emails (
				org_id     varchar(64) not null,
				channel_id varchar(16) not null,
				email      varchar(16) not null,
				primary key(channel_id, email)
			)
			`,
		); err != nil {
			return fmt.Errorf("create tables channel_email: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists channel_user_ids (
				org_id     varchar(64) not null,
				channel_id varchar(16) not null,
				user_id    varchar(16) not null,
				primary key(channel_id, user_id)
			)
			`,
		); err != nil {
			return fmt.Errorf("create tables channel_user_ids: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

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
