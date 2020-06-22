package database

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type APIKeyRepository struct {
	DB *gorm.DB
}

type APIKey struct {
	OrgID      string `gorm:"column:org_id;      type:varchar(16); not null"`
	Name       string `gorm:"column:name;        type:varchar(16); not null"`
	APIKey     string `gorm:"column:api_key;     type:varchar(45); not null; primary key"`
	Read       bool   `gorm:"column:xread;       type:varchar(16); not null; default:1"`
	Write      bool   `gorm:"column:xwrite;      type:varchar(16); not null; default:1"`
	LastAccess int64  `gorm:"column:last_access; type:varchar(16); not null; default:0"`
}

func (k APIKey) Domain() domain.APIKey {
	return domain.APIKey{
		OrgID:      k.OrgID,
		Name:       k.Name,
		APIKey:     k.APIKey,
		Read:       k.Read,
		Write:      k.Write,
		LastAccess: k.LastAccess,
	}
}

func NewAPIKeyRepository(handler SQLHandler) *APIKeyRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

	if err := db.AutoMigrate(&APIKey{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate apikey: %v", err))
	}

	return &APIKeyRepository{
		DB: db,
	}
}

func (repo *APIKeyRepository) Save(orgID, name, apikey string, write bool) (*domain.APIKey, error) {
	create := APIKey{
		OrgID:      orgID,
		Name:       name,
		APIKey:     apikey,
		Read:       true,
		Write:      write,
		LastAccess: time.Now().Unix(),
	}

	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if found := !tx.Where(&APIKey{APIKey: apikey}).First(&APIKey{}).RecordNotFound(); found {
			return nil
		}

		if err := tx.Create(&create).Error; err != nil {
			return fmt.Errorf("create: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	out := create.Domain()
	return &out, nil
}

func (repo *APIKeyRepository) APIKey(apikey string) (*domain.APIKey, error) {
	if apikey == "" {
		return nil, fmt.Errorf("apikey is empty")
	}

	result := APIKey{}
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where(&APIKey{APIKey: apikey}).First(&result).RecordNotFound() {
			return fmt.Errorf("apikey not found")
		}

		now := time.Now().Unix()
		if err := tx.Model(&APIKey{}).Where(&APIKey{APIKey: apikey}).Update(&APIKey{LastAccess: now}).Error; err != nil {
			return fmt.Errorf("update: %v", err)
		}
		result.LastAccess = now

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	out := result.Domain()
	return &out, nil
}
