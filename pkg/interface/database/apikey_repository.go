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

	apikey := APIKey{
		OrgID:      "4b825dc642c",
		Name:       "default",
		APIKey:     "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb",
		Read:       true,
		Write:      true,
		LastAccess: time.Now().Unix(),
	}

	if err := db.Where(&apikey).Assign(&apikey).FirstOrCreate(&APIKey{}).Error; err != nil {
		panic(fmt.Errorf("first or create: %v", err))
	}

	return &APIKeyRepository{
		DB: db,
	}
}

func (repo *APIKeyRepository) APIKey(apikey string) (*domain.APIKey, error) {
	if len(apikey) < 1 {
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
