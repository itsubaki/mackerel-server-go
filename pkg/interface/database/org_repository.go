package database

import (
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type OrgRepository struct {
	DB *gorm.DB
}

type Org struct {
	ID   string `gorm:"column:id;   type:varchar(16); not null; primary_key"`
	Name string `gorm:"column:name; type:varchar(16); not null; unique"`
}

func (o Org) Domain() domain.Org {
	return domain.Org{
		ID:   o.ID,
		Name: o.Name,
	}
}

func NewOrgRepository(handler SQLHandler) *OrgRepository {
	db, err := gorm.Open(mysql.Open(handler.DSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugging() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&Org{}); err != nil {
		panic(fmt.Errorf("auto migrate org: %v", err))
	}

	return &OrgRepository{
		DB: db,
	}
}

func (repo *OrgRepository) Save(orgID, name string) (*domain.Org, error) {
	create := Org{
		ID:   orgID,
		Name: name,
	}

	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&Org{ID: orgID}).First(&Org{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
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

func (repo *OrgRepository) Org(orgID string) (*domain.Org, error) {
	result := Org{}
	if err := repo.DB.Where(&Org{ID: orgID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from orgs: %v", err)
	}

	return &domain.Org{ID: orgID, Name: result.Name}, nil
}
