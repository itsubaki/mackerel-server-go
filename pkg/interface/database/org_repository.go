package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type OrgRepository struct {
	DB *gorm.DB
}

type Org struct {
	ID   string `gorm:"column:id;   type:varchar(16); not null; primary_key"`
	Name string `gorm:"column:name; type:varchar(16); not null; unique"`
}

func NewOrgRepository(handler SQLHandler) *OrgRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())
	if err := db.AutoMigrate(&Org{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate org: %v", err))
	}

	org := Org{ID: "4b825dc642c", Name: "mackerel"}
	if err := db.Where(&org).Assign(&org).FirstOrCreate(&org).Error; err != nil {
		panic(fmt.Errorf("first or create: %v", err))
	}

	return &OrgRepository{
		DB: db,
	}
}

func (repo *OrgRepository) Org(orgID string) (*domain.Org, error) {
	result := Org{}
	if err := repo.DB.Where(&Org{ID: orgID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from orgs: %v", err)
	}

	return &domain.Org{ID: orgID, Name: result.Name}, nil
}
