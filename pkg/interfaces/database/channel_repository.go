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
				type                enum('email', 'slack', 'webhook') not null,
				url                 text,
				enabled_graph_image bool not null default '1'
			)
			`,
		); err != nil {
			return fmt.Errorf("create tables channels: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists channel_mentions (
				org_id     varchar(64) not null,
				channel_id varchar(16) not null,
				status     enum('ok', 'warning', 'critical') not null,
				message    text,
				primary key(channel_id, status),
				foreign key fk_group(channel_id) references channels(id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create tables channel_mentions: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists channel_events (
				org_id     varchar(64) not null,
				channel_id varchar(16) not null,
				event      varchar(16) not null,
				primary key(channel_id, event),
				foreign key fk_group(channel_id) references channels(id) on delete cascade on update cascade
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
				primary key(channel_id, email),
				foreign key fk_group(channel_id) references channels(id) on delete cascade on update cascade
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
				primary key(channel_id, user_id),
				foreign key fk_group(channel_id) references channels(id) on delete cascade on update cascade
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

func (repo *ChannelRepository) mentions(tx Tx, orgID, channelID string) (map[string]string, error) {
	rows, err := tx.Query(`select status, message from channel_mentions where org_id=? and channel_id=?`, orgID, channelID)
	if err != nil {
		return nil, fmt.Errorf("select * from channel_mentions: %v", err)
	}
	defer rows.Close()

	mentions := make(map[string]string)
	for rows.Next() {
		var status, message string
		if err := rows.Scan(
			&status,
			&message,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		mentions[status] = message
	}

	return mentions, nil
}

func (repo *ChannelRepository) events(tx Tx, orgID, channelID string) ([]string, error) {
	rows, err := tx.Query(`select event from channel_events where org_id=? and channel_id=?`, orgID, channelID)
	if err != nil {
		return nil, fmt.Errorf("select * from channel_events: %v", err)
	}
	defer rows.Close()

	events := make([]string, 0)
	for rows.Next() {
		var event string
		if err := rows.Scan(
			&event,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (repo *ChannelRepository) emails(tx Tx, orgID, channelID string) ([]string, error) {
	rows, err := tx.Query(`select email from channel_emails where org_id=? and channel_id=?`, orgID, channelID)
	if err != nil {
		return nil, fmt.Errorf("select * from channel_emails: %v", err)
	}
	defer rows.Close()

	emails := make([]string, 0)
	for rows.Next() {
		var email string
		if err := rows.Scan(
			&email,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func (repo *ChannelRepository) userIDs(tx Tx, orgID, channelID string) ([]string, error) {
	rows, err := tx.Query(`select user_id from channel_user_ids where org_id=? and channel_id=?`, orgID, channelID)
	if err != nil {
		return nil, fmt.Errorf("select * from channel_user_ids: %v", err)
	}
	defer rows.Close()

	userIDs := make([]string, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(
			&userID,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (repo *ChannelRepository) List(orgID string) (*domain.Channels, error) {
	channels := make([]domain.Channel, 0)
	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query(`select * from channels where org_id=?`, orgID)
		if err != nil {
			return fmt.Errorf("selet * from channels: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var channel domain.Channel
			var trash string
			if err := rows.Scan(
				&trash,
				&channel.ID,
				&channel.Name,
				&channel.Type,
				&channel.URL,
				&channel.EnabledGraphImage,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			channels = append(channels, channel)
		}

		for i := range channels {
			mentions, err := repo.mentions(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("mentions: %v", err)
			}
			channels[i].Mentions = mentions

			events, err := repo.events(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("events: %v", err)
			}
			channels[i].Events = events

			emails, err := repo.emails(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("emails: %v", err)
			}
			channels[i].Emails = emails

			userIDs, err := repo.userIDs(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("userIDs: %v", err)
			}
			channels[i].UserIDs = userIDs
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	chanif := make([]interface{}, 0)
	for i := range channels {
		chanif = append(chanif, channels[i].Cast())
	}

	return &domain.Channels{Channels: chanif}, nil
}

func (repo *ChannelRepository) Save(orgID string, channel *domain.Channel) (interface{}, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into channels (
				org_id,
				id,
				name,
				type,
				url,
				enabled_graph_image
			) values (?, ?, ?, ?, ?, ?)
			`,
			orgID,
			channel.ID,
			channel.Name,
			channel.Type,
			channel.URL,
			channel.EnabledGraphImage,
		); err != nil {
			return fmt.Errorf("insert into channels: %v", err)
		}

		for status, message := range channel.Mentions {
			if _, err := tx.Exec(
				`
				insert into channel_mentions (
					org_id,
					channel_id,
					status,
					message
				) values (?, ?, ?, ?)
				`,
				orgID,
				channel.ID,
				status,
				message,
			); err != nil {
				return fmt.Errorf("insert into channel_mentions: %v", err)
			}
		}

		for i := range channel.Events {
			if _, err := tx.Exec(
				`
				insert into channel_events (
					org_id,
					channel_id,
					event,
				) values (?, ?, ?)
				`,
				orgID,
				channel.ID,
				channel.Events[i],
			); err != nil {
				return fmt.Errorf("insert into channel_events: %v", err)
			}
		}

		for i := range channel.Emails {
			if _, err := tx.Exec(
				`
				insert into channel_emails (
					org_id,
					channel_id,
					email,
				) values (?, ?, ?)
				`,
				orgID,
				channel.ID,
				channel.Emails[i],
			); err != nil {
				return fmt.Errorf("insert into channel_emails: %v", err)
			}
		}

		for i := range channel.UserIDs {
			if _, err := tx.Exec(
				`
				insert into channel_user_ids (
					org_id,
					channel_id,
					user_id,
				) values (?, ?, ?)
				`,
				orgID,
				channel.ID,
				channel.UserIDs[i],
			); err != nil {
				return fmt.Errorf("insert into channel_user_ids: %v", err)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return channel.Cast(), nil
}

func (repo *ChannelRepository) Exists(orgID, channelID string) bool {
	rows, err := repo.Query("select 1 from channels where org_id=? and id=?", orgID, channelID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *ChannelRepository) Delete(orgID, channelID string) (interface{}, error) {
	var channel domain.Channel
	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow("select * from channels where org_id=? and id=?", orgID, channelID)

		var trash string
		if err := row.Scan(
			&trash,
			&channel.ID,
			&channel.Name,
			&channel.Type,
			&channel.URL,
			&channel.EnabledGraphImage,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		mentions, err := repo.mentions(tx, orgID, channel.ID)
		if err != nil {
			return fmt.Errorf("mentions: %v", err)
		}
		channel.Mentions = mentions

		events, err := repo.events(tx, orgID, channel.ID)
		if err != nil {
			return fmt.Errorf("events: %v", err)
		}
		channel.Events = events

		emails, err := repo.emails(tx, orgID, channel.ID)
		if err != nil {
			return fmt.Errorf("emails: %v", err)
		}
		channel.Emails = emails

		userIDs, err := repo.userIDs(tx, orgID, channel.ID)
		if err != nil {
			return fmt.Errorf("userIDs: %v", err)
		}
		channel.UserIDs = userIDs

		if _, err := tx.Exec("delete from channels where org_id=? and id=?", orgID, channelID); err != nil {
			return fmt.Errorf("delete from channels: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return channel.Cast(), nil
}
