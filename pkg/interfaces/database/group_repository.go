package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type NotificationGroupRepository struct {
	SQLHandler
}

func NewNotificationGroupRepository(handler SQLHandler) *NotificationGroupRepository {
	if err := handler.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			create table if not exists notification_groups (
				org_id varchar(16) not null,
				id     varchar(16) not null primary key,
				name   varchar(128) not null,
				level  enum('all', 'critical') not null default 'all'
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_groups: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_children (
				org_id   varchar(16) not null,
				group_id varchar(16) not null,
				child_id varchar(16) not null,
				primary key(group_id, child_id),
				foreign key fk_group(group_id) references notification_groups(id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_group_children: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_channels (
				org_id     varchar(16) not null,
				group_id   varchar(16) not null,
				channel_id varchar(16) not null,
				primary key(group_id, channel_id),
				foreign key fk_group(group_id) references notification_groups(id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_group_channels: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_monitors (
				org_id       varchar(16) not null,
				group_id     varchar(16) not null,
				monitor_id   varchar(16) not null,
				skip_default boolean     not null default '0',
				primary key(group_id, monitor_id),
				foreign key fk_group(group_id) references notification_groups(id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_group_monitors: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_services (
				org_id       varchar(16) not null,
				group_id     varchar(16) not null,
				service_name varchar(128) not null,
				primary key(group_id, service_name),
				foreign key fk_group(group_id) references notification_groups(id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_group_servivces: %v", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("transaction: %v", err))
	}

	return &NotificationGroupRepository{
		SQLHandler: handler,
	}
}

func (repo *NotificationGroupRepository) groupIDs(tx Tx, orgID, groupID string) ([]string, error) {
	rows, err := tx.Query(`select child_id from notification_group_children where org_id=? and group_id=?`, orgID, groupID)
	if err != nil {
		return nil, fmt.Errorf("select * from notification_group_children: %v", err)
	}
	defer rows.Close()

	children := make([]string, 0)
	for rows.Next() {
		var childID string
		if err := rows.Scan(
			&childID,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		children = append(children, childID)
	}

	return children, nil
}

func (repo *NotificationGroupRepository) channelIDs(tx Tx, orgID, groupID string) ([]string, error) {
	rows, err := tx.Query(`select channel_id from notification_group_channels where org_id=? and group_id=?`, orgID, groupID)
	if err != nil {
		return nil, fmt.Errorf("select * from notification_group_channels: %v", err)
	}
	defer rows.Close()

	channels := make([]string, 0)
	for rows.Next() {
		var channelID string
		if err := rows.Scan(
			&channelID,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		channels = append(channels, channelID)
	}

	return channels, nil
}

func (repo *NotificationGroupRepository) monitors(tx Tx, orgID, groupID string) ([]domain.NMonitor, error) {
	rows, err := tx.Query(`select monitor_id, skip_default from notification_group_monitors where org_id=? and group_id=?`, orgID, groupID)
	if err != nil {
		return nil, fmt.Errorf("select * from notification_group_monitors: %v", err)
	}
	defer rows.Close()

	monitors := make([]domain.NMonitor, 0)
	for rows.Next() {
		var monitor domain.NMonitor
		if err := rows.Scan(
			&monitor.ID,
			&monitor.SkipDefault,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

func (repo *NotificationGroupRepository) services(tx Tx, orgID, groupID string) ([]domain.NService, error) {
	rows, err := tx.Query(`select service_name from notification_group_services where org_id=? and group_id=?`, orgID, groupID)
	if err != nil {
		return nil, fmt.Errorf("select * from notification_group_services: %v", err)
	}
	defer rows.Close()

	services := make([]domain.NService, 0)
	for rows.Next() {
		var service domain.NService
		if err := rows.Scan(
			&service.Name,
		); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		services = append(services, service)
	}

	return services, nil
}

func (repo *NotificationGroupRepository) List(orgID string) (*domain.NotificationGroups, error) {
	groups := make([]domain.NotificationGroup, 0)

	if err := repo.Transact(func(tx Tx) error {
		rows, err := tx.Query(`select * from notification_groups where org_id=?`, orgID)
		if err != nil {
			return fmt.Errorf("selet * from notification_groups: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var group domain.NotificationGroup
			if err := rows.Scan(
				&group.OrgID,
				&group.ID,
				&group.Name,
				&group.NotificationLevel,
			); err != nil {
				return fmt.Errorf("scan: %v", err)
			}

			groups = append(groups, group)
		}

		for i := range groups {
			groupIDs, err := repo.groupIDs(tx, orgID, groups[i].ID)
			if err != nil {
				return fmt.Errorf("group_id: %v", err)
			}
			groups[i].ChildNotificationGroupIDs = groupIDs

			channelIDs, err := repo.channelIDs(tx, orgID, groups[i].ID)
			if err != nil {
				return fmt.Errorf("channel_id: %v", err)
			}
			groups[i].ChildChannelIDs = channelIDs

			monitors, err := repo.monitors(tx, orgID, groups[i].ID)
			if err != nil {
				return fmt.Errorf("monitors: %v", err)
			}
			groups[i].Monitors = monitors

			services, err := repo.services(tx, orgID, groups[i].ID)
			if err != nil {
				return fmt.Errorf("services: %v", err)
			}
			groups[i].Services = services
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.NotificationGroups{NotificationGroups: groups}, nil
}

func (repo *NotificationGroupRepository) Exists(orgID, groupID string) bool {
	rows, err := repo.Query("select 1 from notification_groups where org_id=? and id=?", orgID, groupID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func (repo *NotificationGroupRepository) Save(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec(
			`
			insert into notification_groups (
				org_id,
				id,
				name,
				level
			) values (?, ?, ?, ?)
			`,
			orgID,
			group.ID,
			group.Name,
			group.NotificationLevel,
		); err != nil {
			return fmt.Errorf("insert into notification_groups: %v", err)
		}

		for i := range group.ChildNotificationGroupIDs {
			if _, err := tx.Exec(
				`
			insert into notification_group_children (
				org_id,
				group_id,
				child_id
			) values (?, ?, ?)
			`,
				orgID,
				group.ID,
				group.ChildNotificationGroupIDs[i],
			); err != nil {
				return fmt.Errorf("insert into notification_group_children: %v", err)
			}
		}

		for i := range group.ChildChannelIDs {
			if _, err := tx.Exec(
				`
			insert into notification_group_channels (
				org_id,
				group_id,
				channel_id
			) values (?, ?, ?)
			`,
				orgID,
				group.ID,
				group.ChildChannelIDs[i],
			); err != nil {
				return fmt.Errorf("insert into notification_group_channels: %v", err)
			}
		}

		for i := range group.Monitors {
			if _, err := tx.Exec(
				`
			insert into notification_group_monitors (
				org_id,
				group_id,
				monitor_id,
				skip_default
			) values (?, ?, ?, ?)
			`,
				orgID,
				group.ID,
				group.Monitors[i].ID,
				group.Monitors[i].SkipDefault,
			); err != nil {
				return fmt.Errorf("insert into notification_group_monitors: %v", err)
			}
		}

		for i := range group.Services {
			if _, err := tx.Exec(
				`
			insert into notification_group_services (
				org_id,
				group_id,
				service_name
			) values (?, ?, ?)
			`,
				orgID,
				group.ID,
				group.Services[i].Name,
			); err != nil {
				return fmt.Errorf("insert into notification_group_services: %v", err)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return group, nil
}

func (repo *NotificationGroupRepository) Update(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	if err := repo.Transact(func(tx Tx) error {
		if _, err := tx.Exec("delete from notification_groups where org_id=? and id=?", orgID, group.ID); err != nil {
			return fmt.Errorf("delete from notification_groups: %v", err)
		}

		if _, err := tx.Exec(
			`
			insert into notification_groups (
				org_id,
				id,
				name,
				level
			) values (?, ?, ?, ?)
			`,
			orgID,
			group.ID,
			group.Name,
			group.NotificationLevel,
		); err != nil {
			return fmt.Errorf("insert into notification_groups: %v", err)
		}

		for i := range group.ChildNotificationGroupIDs {
			if _, err := tx.Exec(
				`
			insert into notification_group_children (
				org_id,
				group_id,
				child_id
			) values (?, ?, ?)
			`,
				orgID,
				group.ID,
				group.ChildNotificationGroupIDs[i],
			); err != nil {
				return fmt.Errorf("insert into notification_group_children: %v", err)
			}
		}

		for i := range group.ChildChannelIDs {
			if _, err := tx.Exec(
				`
			insert into notification_group_channels (
				org_id,
				group_id,
				channel_id
			) values (?, ?, ?)
			`,
				orgID,
				group.ID,
				group.ChildChannelIDs[i],
			); err != nil {
				return fmt.Errorf("insert into notification_group_channels: %v", err)
			}
		}

		for i := range group.Monitors {
			if _, err := tx.Exec(
				`
			insert into notification_group_monitors (
				org_id,
				group_id,
				monitor_id,
				skip_default
			) values (?, ?, ?, ?)
			`,
				orgID,
				group.ID,
				group.Monitors[i].ID,
				group.Monitors[i].SkipDefault,
			); err != nil {
				return fmt.Errorf("insert into notification_group_monitors: %v", err)
			}
		}

		for i := range group.Services {
			if _, err := tx.Exec(
				`
			insert into notification_group_services (
				org_id,
				group_id,
				service_name
			) values (?, ?, ?)
			`,
				orgID,
				group.ID,
				group.Services[i].Name,
			); err != nil {
				return fmt.Errorf("insert into notification_group_services: %v", err)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return group, nil
}

func (repo *NotificationGroupRepository) Delete(orgID, groupID string) (*domain.NotificationGroup, error) {
	group := domain.NotificationGroup{}

	if err := repo.Transact(func(tx Tx) error {
		row := tx.QueryRow(`select * from notification_groups where org_id=? and id=?`, orgID, groupID)
		if err := row.Scan(
			&group.OrgID,
			&group.ID,
			&group.Name,
			&group.NotificationLevel,
		); err != nil {
			return fmt.Errorf("scan: %v", err)
		}

		groupIDs, err := repo.groupIDs(tx, orgID, groupID)
		if err != nil {
			return fmt.Errorf("group_id: %v", err)
		}
		group.ChildNotificationGroupIDs = groupIDs

		channelIDs, err := repo.channelIDs(tx, orgID, groupID)
		if err != nil {
			return fmt.Errorf("channel_id: %v", err)
		}
		group.ChildChannelIDs = channelIDs

		monitors, err := repo.monitors(tx, orgID, groupID)
		if err != nil {
			return fmt.Errorf("monitors: %v", err)
		}
		group.Monitors = monitors

		services, err := repo.services(tx, orgID, groupID)
		if err != nil {
			return fmt.Errorf("services: %v", err)
		}
		group.Services = services

		if _, err := tx.Exec("delete from notification_groups where org_id=? and id=?", orgID, groupID); err != nil {
			return fmt.Errorf("delete from notification_groups: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &group, nil
}
