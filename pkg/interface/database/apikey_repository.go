package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.APIKeyRepository = (*APIKeyRepository)(nil)

type APIKeyRepository struct {
	DB *gorm.DB
}

type APIKey struct {
	OrgID      string `gorm:"column:org_id;      type:varchar(16); not null"`
	Name       string `gorm:"column:name;        type:varchar(16); not null"`
	APIKey     string `gorm:"column:api_key;     type:varchar(45); not null; primary_key"`
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
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&APIKey{}); err != nil {
		panic(fmt.Errorf("auto migrate apikey: %v", err))
	}

	return &APIKeyRepository{
		DB: db,
	}
}

func (r *APIKeyRepository) Save(orgID, name, apikey string, write bool) (*domain.APIKey, error) {
	create := APIKey{
		OrgID:      orgID,
		Name:       name,
		APIKey:     apikey,
		Read:       true,
		Write:      write,
		LastAccess: time.Now().Unix(),
	}

	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&APIKey{}).Where(&APIKey{APIKey: apikey}).Count(&count).Error; err != nil {
			return fmt.Errorf("count: %v", err)
		}

		if count == 0 {
			if err := tx.Create(&create).Error; err != nil {
				return fmt.Errorf("create: %v", err)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction: %v", err)
	}

	out := create.Domain()
	return &out, nil
}

func (r *APIKeyRepository) APIKey(apikey string) (*domain.APIKey, error) {
	if apikey == "" {
		return nil, fmt.Errorf("apikey is empty")
	}

	result := APIKey{}
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&APIKey{APIKey: apikey}).First(&result).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("apikey not found")
		}

		now := time.Now().Unix()
		if err := tx.Model(&APIKey{}).Where(&APIKey{APIKey: apikey}).Updates(&APIKey{LastAccess: now}).Error; err != nil {
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
