package database

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	DB *gorm.DB
}

type Channel struct {
	OrgID             string `gorm:"column:org_id;              type:varchar(16);  not null"`
	ID                string `gorm:"column:id;                  type:varchar(16);  not null; primary_key"`
	Name              string `gorm:"column:name;                type:varchar(16);  not null"`
	Type              string `gorm:"column:type;                type:enum('email', 'slack', 'webhook');  not null"`
	URL               string `gorm:"column:url;                 type:text;"`
	EnabledGraphImage bool   `gorm:"column:enabled_graph_image; type:bool; not null; default:true"`
}

type ChannelMention struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16);  not null"`
	ChannelID string `gorm:"column:channel_id; type:varchar(16);  not null; primary_key"`
	Status    string `gorm:"column:status;     type:enum('ok', 'warning', 'critical');  not null; primary_key"`
	Message   string `gorm:"column:message;    type:text;"`
}

type ChannelEvent struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16);  not null"`
	ChannelID string `gorm:"column:channel_id; type:varchar(16);  not null; primary_key"`
	Event     string `gorm:"column:event;      type:varchar(16);  not null; primary_key"`
}

type ChannelEmail struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16);  not null"`
	ChannelID string `gorm:"column:channel_id; type:varchar(16);  not null; primary_key"`
	EMail     string `gorm:"column:email;      type:varchar(128); not null; primary_key"`
}

type ChannelUserID struct {
	OrgID     string `gorm:"column:org_id;     type:varchar(16); not null"`
	ChannelID string `gorm:"column:channel_id; type:varchar(16); not null; primary_key"`
	UserID    string `gorm:"column:user_id;    type:varchar(16); not null; primary_key"`
}

func NewChannelRepository(handler SQLHandler) *ChannelRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&Channel{}); err != nil {
		panic(fmt.Errorf("auto migrate channel: %v", err))
	}

	if err := db.AutoMigrate(&ChannelMention{}); err != nil {
		panic(fmt.Errorf("auto migrate channel mention: %v", err))
	}

	if err := db.AutoMigrate(&ChannelEvent{}); err != nil {
		panic(fmt.Errorf("auto migrate channel event: %v", err))
	}

	if err := db.AutoMigrate(&ChannelEmail{}); err != nil {
		panic(fmt.Errorf("auto migrate channel email: %v", err))
	}

	if err := db.AutoMigrate(&ChannelUserID{}); err != nil {
		panic(fmt.Errorf("auto migrate channel user id: %v", err))
	}

	return &ChannelRepository{
		DB: db,
	}
}

func (r *ChannelRepository) List(orgID string) (*domain.Channels, error) {
	channels := make([]domain.Channel, 0)
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := make([]Channel, 0)
		if err := tx.Where(&Channel{OrgID: orgID}).Find(&result).Error; err != nil {
			return fmt.Errorf("selet * from channels: %v", err)
		}

		for _, r := range result {
			channels = append(channels, domain.Channel{
				OrgID:             r.OrgID,
				ID:                r.ID,
				Name:              r.Name,
				Type:              r.Type,
				URL:               r.URL,
				EnabledGraphImage: r.EnabledGraphImage,
			})
		}

		for i := range channels {
			mentions, err := r.mentions(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("mentions: %v", err)
			}
			channels[i].Mentions = mentions

			events, err := r.events(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("events: %v", err)
			}
			channels[i].Events = events

			emails, err := r.emails(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("emails: %v", err)
			}
			channels[i].Emails = emails

			userIDs, err := r.userIDs(tx, orgID, channels[i].ID)
			if err != nil {
				return fmt.Errorf("userIDs: %v", err)
			}
			channels[i].UserIDs = userIDs
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	out := make([]interface{}, 0)
	for i := range channels {
		out = append(out, channels[i].Cast())
	}

	return &domain.Channels{Channels: out}, nil
}

func (r *ChannelRepository) Save(orgID string, channel *domain.Channel) (interface{}, error) {
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Channel{
			OrgID:             orgID,
			ID:                channel.ID,
			Name:              channel.Name,
			Type:              channel.Type,
			URL:               channel.URL,
			EnabledGraphImage: channel.EnabledGraphImage,
		}).Error; err != nil {
			return fmt.Errorf("insert into channels: %v", err)
		}

		for st, mes := range channel.Mentions {
			if err := tx.Create(&ChannelMention{
				OrgID:     orgID,
				ChannelID: channel.ID,
				Status:    st,
				Message:   mes,
			}).Error; err != nil {
				return fmt.Errorf("insert into channel_mentions: %v", err)
			}
		}

		for i := range channel.Events {
			if err := tx.Create(&ChannelEvent{
				OrgID:     orgID,
				ChannelID: channel.ID,
				Event:     channel.Events[i],
			}).Error; err != nil {
				return fmt.Errorf("insert into channel_events: %v", err)
			}
		}

		for i := range channel.Emails {
			if err := tx.Create(&ChannelEmail{
				OrgID:     orgID,
				ChannelID: channel.ID,
				EMail:     channel.Emails[i],
			}).Error; err != nil {
				return fmt.Errorf("insert into channel_emails: %v", err)
			}
		}

		for i := range channel.UserIDs {
			if err := tx.Create(&ChannelUserID{
				OrgID:     orgID,
				ChannelID: channel.ID,
				UserID:    channel.UserIDs[i],
			}).Error; err != nil {
				return fmt.Errorf("insert into channel_emails: %v", err)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return channel.Cast(), nil
}

func (r *ChannelRepository) Exists(orgID, channelID string) bool {
	if err := r.DB.Where(&Channel{OrgID: orgID, ID: channelID}).First(&Channel{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (r *ChannelRepository) Delete(orgID, channelID string) (interface{}, error) {
	var channel domain.Channel
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := Channel{}
		if err := tx.Where(&Channel{OrgID: orgID, ID: channelID}).Find(&result).Error; err != nil {
			return fmt.Errorf("select * from channels: %v", err)
		}

		channel.OrgID = result.OrgID
		channel.ID = result.ID
		channel.Name = result.Name
		channel.Type = result.Type
		channel.URL = result.URL
		channel.EnabledGraphImage = result.EnabledGraphImage

		mentions, err := r.mentions(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("mentions: %v", err)
		}
		channel.Mentions = mentions

		events, err := r.events(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("events: %v", err)
		}
		channel.Events = events

		emails, err := r.emails(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("emails: %v", err)
		}
		channel.Emails = emails

		userIDs, err := r.userIDs(tx, orgID, result.ID)
		if err != nil {
			return fmt.Errorf("userIDs: %v", err)
		}
		channel.UserIDs = userIDs

		if err := tx.Delete(&Channel{OrgID: orgID, ID: channelID}).Error; err != nil {
			return fmt.Errorf("delete from channels: %v", err)
		}

		if err := tx.Where("org_id = ? AND channel_id = ?", orgID, channelID).Delete(&ChannelMention{}).Error; err != nil {
			return fmt.Errorf("delete from channel_mentions: %v", err)
		}

		if err := tx.Where("org_id = ? AND channel_id = ?", orgID, channelID).Delete(&ChannelEvent{}).Error; err != nil {
			return fmt.Errorf("delete from channel_events: %v", err)
		}

		if err := tx.Where("org_id = ? AND channel_id = ?", orgID, channelID).Delete(&ChannelEmail{}).Error; err != nil {
			return fmt.Errorf("delete from channel_emails: %v", err)
		}

		if err := tx.Where("org_id = ? AND channel_id = ?", orgID, channelID).Delete(&ChannelUserID{}).Error; err != nil {
			return fmt.Errorf("delete from channel_user_ids: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	return channel.Cast(), nil
}

func (r *ChannelRepository) mentions(tx *gorm.DB, orgID, channelID string) (map[string]string, error) {
	result := make([]ChannelMention, 0)
	if err := tx.Where(&ChannelMention{OrgID: orgID, ChannelID: channelID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from channel_mentions: %v", err)
	}

	mentions := make(map[string]string)
	for _, r := range result {
		mentions[r.Status] = r.Message
	}

	return mentions, nil
}

func (r *ChannelRepository) events(tx *gorm.DB, orgID, channelID string) ([]string, error) {
	result := make([]ChannelEvent, 0)
	if err := tx.Where(&ChannelEvent{OrgID: orgID, ChannelID: channelID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from channel_events: %v", err)
	}

	events := make([]string, 0)
	for _, r := range result {
		events = append(events, r.Event)
	}

	return events, nil
}

func (r *ChannelRepository) emails(tx *gorm.DB, orgID, channelID string) ([]string, error) {
	result := make([]ChannelEmail, 0)
	if err := tx.Where(&ChannelEmail{OrgID: orgID, ChannelID: channelID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from channel_emails: %v", err)
	}

	emails := make([]string, 0)
	for _, r := range result {
		emails = append(emails, r.EMail)
	}

	return emails, nil
}

func (r *ChannelRepository) userIDs(tx *gorm.DB, orgID, channelID string) ([]string, error) {
	result := make([]ChannelUserID, 0)
	if err := tx.Where(&ChannelUserID{OrgID: orgID, ChannelID: channelID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from channel_user_ids: %v", err)
	}

	userIDs := make([]string, 0)
	for _, r := range result {
		userIDs = append(userIDs, r.UserID)
	}

	return userIDs, nil
}
