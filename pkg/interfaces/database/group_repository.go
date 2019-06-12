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
				org_id varchar(64) not null,
				id     varchar(16) not null primary key,
				name   varchar(16) not null,
				level  enum('all', 'critical') not null default 'all'
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_groups: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_children (
				org_id   varchar(64) not null,
				group_id varchar(16) not null,
				child_id varchar(16) not null,
				primary key(group_id, child_id),
				foreign key fk_group(org_id, group_id) references notification_groups(org_id, id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_group_children: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_channels (
				org_id     varchar(64) not null,
				group_id   varchar(16) not null,
				channel_id varchar(16) not null,
				primary key(group_id, channel_id),
				foreign key fk_group(org_id, group_id) references notification_groups(org_id, id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_group_channels: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_monitors (
				org_id     varchar(64) not null,
				group_id   varchar(16) not null,
				monitor_id varchar(16) not null,
				primary key(group_id, monitor_id),
				foreign key fk_group(org_id, group_id) references notification_groups(org_id, id) on delete cascade on update cascade
			)
			`,
		); err != nil {
			return fmt.Errorf("create table notification_group_monitors: %v", err)
		}

		if _, err := tx.Exec(
			`
			create table if not exists notification_group_services (
				org_id       varchar(64) not null,
				group_id     varchar(16) not null,
				service_name varchar(16) not null,
				primary key(group_id, service_name),
				foreign key fk_group(org_id, group_id) references notification_groups(org_id, id) on delete cascade on update cascade
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

func (repo *NotificationGroupRepository) List(orgID string) (*domain.NotificationGroups, error) {
	return nil, nil
}

func (repo *NotificationGroupRepository) Save(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	return nil, nil
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

func (repo *NotificationGroupRepository) Update(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	return nil, nil
}

func (repo *NotificationGroupRepository) Delete(orgID, groupID string) (*domain.NotificationGroup, error) {
	return nil, nil
}
