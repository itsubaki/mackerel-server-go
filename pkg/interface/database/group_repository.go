package database

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type NotificationGroupRepository struct {
	DB *gorm.DB
}

type NotificationGroup struct {
	OrgID             string `gorm:"column:org_id; type:varchar(16);  not null;"`
	ID                string `gorm:"column:id;     type:varchar(16);  not null; primary_key"`
	Name              string `gorm:"column:name;   type:varchar(128); not null;"`
	NotificationLevel string `gorm:"column:level;  type:enum('all', 'critical'); not null; default:all"`
}

type NotificationGroupChild struct {
	OrgID   string `gorm:"column:org_id;   type:varchar(16);  not null;"`
	GroupID string `gorm:"column:group_id; type:varchar(16);  not null; primary_key"`
	ChildID string `gorm:"column:child_id; type:varchar(128); not null; primary_key"`
}

type NotificationGroupChannel struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16); not null;"`
	GroupID   string `gorm:"column:group_id;   type:varchar(16); not null; primary_key"`
	ChannelID string `gorm:"column:channel_id; type:varchar(16); not null; primary_key"`
}

type NotificationGroupMonitor struct {
	OrgID       string `gorm:"column:org_id;       type:varchar(16); not null;"`
	GroupID     string `gorm:"column:group_id;     type:varchar(16); not null; primary_key"`
	MonitorID   string `gorm:"column:monitor_id;   type:varchar(16); not null; primary_key"`
	SkipDefault bool   `gorm:"column:skip_default; type:boolean; not null; default:false"`
}

type NotificationGroupService struct {
	OrgID       string `gorm:"column:org_id;       type:varchar(16);  not null;"`
	GroupID     string `gorm:"column:group_id;     type:varchar(16);  not null; primary_key"`
	ServiceName string `gorm:"column:service_name; type:varchar(128); not null; primary_key"`
}

func NewNotificationGroupRepository(handler SQLHandler) *NotificationGroupRepository {
	db, err := gorm.Open(mysql.Open(handler.DSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugging() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&NotificationGroup{}); err != nil {
		panic(fmt.Errorf("auto migrate notification_groups: %v", err))
	}

	if err := db.AutoMigrate(&NotificationGroupChild{}); err != nil {
		panic(fmt.Errorf("auto migrate notification_group_children: %v", err))
	}

	if err := db.AutoMigrate(&NotificationGroupChannel{}); err != nil {
		panic(fmt.Errorf("auto migrate notification_group_channels: %v", err))
	}

	if err := db.AutoMigrate(&NotificationGroupMonitor{}); err != nil {
		panic(fmt.Errorf("auto migrate notification_group_monitors: %v", err))
	}

	if err := db.AutoMigrate(&NotificationGroupService{}); err != nil {
		panic(fmt.Errorf("auto migrate notification_group_services: %v", err))
	}

	return &NotificationGroupRepository{
		DB: db,
	}
}

func (repo *NotificationGroupRepository) groupIDs(tx *gorm.DB, orgID, groupID string) ([]string, error) {
	result := make([]NotificationGroupChild, 0)
	if err := tx.Where(&NotificationGroupChild{OrgID: orgID, GroupID: groupID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from notification_group_children: %v", err)
	}

	out := make([]string, 0)
	for _, r := range result {
		out = append(out, r.ChildID)
	}

	return out, nil
}

func (repo *NotificationGroupRepository) channelIDs(tx *gorm.DB, orgID, groupID string) ([]string, error) {
	result := make([]NotificationGroupChannel, 0)
	if err := tx.Where(&NotificationGroupChannel{OrgID: orgID, GroupID: groupID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from notification_group_children: %v", err)
	}

	out := make([]string, 0)
	for _, r := range result {
		out = append(out, r.ChannelID)
	}

	return out, nil
}

func (repo *NotificationGroupRepository) monitors(tx *gorm.DB, orgID, groupID string) ([]domain.NotificationMonitor, error) {
	result := make([]NotificationGroupMonitor, 0)
	if err := tx.Where(&NotificationGroupMonitor{OrgID: orgID, GroupID: groupID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from notification_group_monitors: %v", err)
	}

	out := make([]domain.NotificationMonitor, 0)
	for _, r := range result {
		out = append(out, domain.NotificationMonitor{
			ID:          r.MonitorID,
			SkipDefault: r.SkipDefault,
		})
	}

	return out, nil
}

func (repo *NotificationGroupRepository) services(tx *gorm.DB, orgID, groupID string) ([]domain.NotificationService, error) {
	result := make([]NotificationGroupService, 0)
	if err := tx.Where(&NotificationGroupService{OrgID: orgID, GroupID: groupID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from notification_group_services: %v", err)
	}

	out := make([]domain.NotificationService, 0)
	for _, r := range result {
		out = append(out, domain.NotificationService{
			Name: r.ServiceName,
		})
	}

	return out, nil
}

func (repo *NotificationGroupRepository) List(orgID string) (*domain.NotificationGroups, error) {
	out := make([]domain.NotificationGroup, 0)
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		result := make([]NotificationGroup, 0)
		if err := tx.Where(&NotificationGroup{OrgID: orgID}).Find(&result).Error; err != nil {
			return fmt.Errorf("selet * from notification_groups: %v", err)
		}

		for _, r := range result {
			out = append(out, domain.NotificationGroup{
				OrgID:             r.OrgID,
				ID:                r.ID,
				Name:              r.Name,
				NotificationLevel: r.NotificationLevel,
			})
		}

		for i := range out {
			groupIDs, err := repo.groupIDs(tx, orgID, out[i].ID)
			if err != nil {
				return fmt.Errorf("group_id: %v", err)
			}
			out[i].ChildNotificationGroupIDs = groupIDs

			channelIDs, err := repo.channelIDs(tx, orgID, out[i].ID)
			if err != nil {
				return fmt.Errorf("channel_id: %v", err)
			}
			out[i].ChildChannelIDs = channelIDs

			monitors, err := repo.monitors(tx, orgID, out[i].ID)
			if err != nil {
				return fmt.Errorf("monitors: %v", err)
			}
			out[i].Monitors = monitors

			services, err := repo.services(tx, orgID, out[i].ID)
			if err != nil {
				return fmt.Errorf("services: %v", err)
			}
			out[i].Services = services
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &domain.NotificationGroups{NotificationGroups: out}, nil
}

func (repo *NotificationGroupRepository) Exists(orgID, groupID string) bool {
	if err := repo.DB.Where(&NotificationGroup{OrgID: orgID, ID: groupID}).Find(&NotificationGroup{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (repo *NotificationGroupRepository) Save(orgID string, group *domain.NotificationGroup) (*domain.NotificationGroup, error) {
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&NotificationGroup{OrgID: orgID, ID: group.ID, Name: group.Name, NotificationLevel: group.NotificationLevel}).Error; err != nil {
			return fmt.Errorf("insert into notification_groups: %v", err)
		}

		for i := range group.ChildNotificationGroupIDs {
			if err := tx.Create(&NotificationGroupChild{OrgID: orgID, GroupID: group.ID, ChildID: group.ChildNotificationGroupIDs[i]}).Error; err != nil {
				return fmt.Errorf("insert into notification_group_children: %v", err)
			}
		}

		for i := range group.ChildChannelIDs {
			if err := tx.Create(&NotificationGroupChannel{OrgID: orgID, GroupID: group.ID, ChannelID: group.ChildChannelIDs[i]}).Error; err != nil {
				return fmt.Errorf("insert into notification_group_channels: %v", err)
			}
		}

		for i := range group.Monitors {
			if err := tx.Create(&NotificationGroupMonitor{OrgID: orgID, GroupID: group.ID, MonitorID: group.Monitors[i].ID, SkipDefault: group.Monitors[i].SkipDefault}).Error; err != nil {
				return fmt.Errorf("insert into notification_group_monitors: %v", err)
			}
		}

		for i := range group.Services {
			if err := tx.Create(&NotificationGroupService{OrgID: orgID, GroupID: group.ID, ServiceName: group.Services[i].Name}).Error; err != nil {
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
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("org_id = ? AND id = ?", orgID, group.ID).Delete(&NotificationGroup{}).Error; err != nil {
			return fmt.Errorf("delete from notification_groups: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, group.ID).Delete(&NotificationGroupChild{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_children: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, group.ID).Delete(&NotificationGroupChannel{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_channels: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, group.ID).Delete(&NotificationGroupMonitor{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_monitors: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, group.ID).Delete(&NotificationGroupService{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_services: %v", err)
		}

		if err := tx.Create(&NotificationGroup{OrgID: orgID, ID: group.ID, Name: group.Name, NotificationLevel: group.NotificationLevel}).Error; err != nil {
			return fmt.Errorf("insert into notification_groups: %v", err)
		}

		for i := range group.ChildNotificationGroupIDs {
			if err := tx.Create(&NotificationGroupChild{OrgID: orgID, GroupID: group.ID, ChildID: group.ChildNotificationGroupIDs[i]}).Error; err != nil {
				return fmt.Errorf("insert into notification_group_children: %v", err)
			}
		}

		for i := range group.ChildChannelIDs {
			if err := tx.Create(&NotificationGroupChannel{OrgID: orgID, GroupID: group.ID, ChannelID: group.ChildChannelIDs[i]}).Error; err != nil {
				return fmt.Errorf("insert into notification_group_channels: %v", err)
			}
		}

		for i := range group.Monitors {
			if err := tx.Create(&NotificationGroupMonitor{OrgID: orgID, GroupID: group.ID, MonitorID: group.Monitors[i].ID, SkipDefault: group.Monitors[i].SkipDefault}).Error; err != nil {
				return fmt.Errorf("insert into notification_group_monitors: %v", err)
			}
		}

		for i := range group.Services {
			if err := tx.Create(&NotificationGroupService{OrgID: orgID, GroupID: group.ID, ServiceName: group.Services[i].Name}).Error; err != nil {
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
	out := domain.NotificationGroup{}
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		result := NotificationGroup{}
		if err := tx.Where(&NotificationGroup{OrgID: orgID, ID: groupID}).Find(&result).Error; err != nil {
			return fmt.Errorf("selet * from notification_groups: %v", err)
		}

		out.OrgID = result.OrgID
		out.ID = result.ID
		out.Name = result.Name
		out.NotificationLevel = result.NotificationLevel

		groupIDs, err := repo.groupIDs(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("group_id: %v", err)
		}
		out.ChildNotificationGroupIDs = groupIDs

		channelIDs, err := repo.channelIDs(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("channel_id: %v", err)
		}
		out.ChildChannelIDs = channelIDs

		monitors, err := repo.monitors(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("monitors: %v", err)
		}
		out.Monitors = monitors

		services, err := repo.services(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("services: %v", err)
		}
		out.Services = services

		if err := tx.Where("org_id = ? AND id = ?", orgID, groupID).Delete(&NotificationGroup{}).Error; err != nil {
			return fmt.Errorf("delete from notification_groups: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, groupID).Delete(&NotificationGroupChild{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_children: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, groupID).Delete(&NotificationGroupChannel{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_channels: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, groupID).Delete(&NotificationGroupMonitor{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_monitors: %v", err)
		}

		if err := tx.Where("org_id = ? AND group_id = ?", orgID, groupID).Delete(&NotificationGroupService{}).Error; err != nil {
			return fmt.Errorf("delete from notification_group_services: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return &out, nil
}
