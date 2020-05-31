package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type ServiceRepository struct {
	DB *gorm.DB
}

type Service struct {
	OrgID string `gorm:"column:org_id; type:varchar(16);  not null; primary_key"`
	Name  string `gorm:"column:name;   type:varchar(128); not null; primary_key"`
	Memo  string `gorm:"column:memo;   type:varchar(218); not null; default:''"`
}

func NewServiceRepository(handler SQLHandler) *ServiceRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())
	if err := db.AutoMigrate(&Service{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate service: %v", err))
	}

	return &ServiceRepository{
		DB: db,
	}
}

func (repo *ServiceRepository) List(orgID string) (*domain.Services, error) {
	result := make([]Service, 0)
	if err := repo.DB.Where(&Service{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from services: %v", err)
	}

	out := make([]domain.Service, 0)
	for _, r := range result {
		out = append(out, domain.Service{
			OrgID: r.OrgID,
			Name:  r.Name,
			Memo:  r.Memo,
			Roles: make([]string, 0),
		})
	}

	return &domain.Services{Services: out}, nil
}

func (repo *ServiceRepository) Save(orgID string, s *domain.Service) error {
	service := Service{
		OrgID: orgID,
		Name:  s.Name,
		Memo:  s.Memo,
	}

	if err := repo.DB.Where(&service).Assign(&service).FirstOrCreate(&service).Error; err != nil {
		panic(fmt.Errorf("first or create: %v", err))
	}

	return nil
}

func (repo *ServiceRepository) Service(orgID, serviceName string) (*domain.Service, error) {
	result := Service{}
	if err := repo.DB.Where(&Service{OrgID: orgID, Name: serviceName}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from serviecs: %v", err)
	}

	service := domain.Service{
		OrgID: result.OrgID,
		Name:  result.Name,
		Memo:  result.Memo,
		Roles: make([]string, 0),
	}

	return &service, nil
}

func (repo *ServiceRepository) Exists(orgID, serviceName string) bool {
	if repo.DB.Where(&Service{OrgID: orgID, Name: serviceName}).First(&Invitation{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *ServiceRepository) Delete(orgID, serviceName string) error {
	if err := repo.DB.Delete(&Service{OrgID: orgID, Name: serviceName}).Error; err != nil {
		return fmt.Errorf("delete from services: %v", err)
	}

	return nil
}
