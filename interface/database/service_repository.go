package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.ServiceRepository = (*ServiceRepository)(nil)

type ServiceRepository struct {
	DB *gorm.DB
}

type Service struct {
	OrgID string `gorm:"column:org_id; type:varchar(16);  not null; primary_key"`
	Name  string `gorm:"column:name;   type:varchar(128); not null; primary_key"`
	Memo  string `gorm:"column:memo;   type:varchar(218); not null; default:''"`
}

func (s Service) Domain() domain.Service {
	return domain.Service{
		OrgID: s.OrgID,
		Name:  s.Name,
		Memo:  s.Memo,
		Roles: make([]string, 0),
	}
}

func NewServiceRepository(handler SQLHandler) *ServiceRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&Service{}); err != nil {
		panic(fmt.Errorf("auto migrate service: %v", err))
	}

	return &ServiceRepository{
		DB: db,
	}
}

func (r *ServiceRepository) List(orgID string) (*domain.Services, error) {
	result := make([]Service, 0)
	if err := r.DB.Where(&Service{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from services: %v", err)
	}

	out := make([]domain.Service, 0)
	for _, r := range result {
		out = append(out, r.Domain())
	}

	return &domain.Services{Services: out}, nil
}

func (r *ServiceRepository) Save(orgID string, s *domain.Service) error {
	if err := r.DB.Where(&Service{OrgID: orgID, Name: s.Name}).Assign(&Service{Memo: s.Memo}).FirstOrCreate(&Service{}).Error; err != nil {
		return fmt.Errorf("first or create: %v", err)
	}

	return nil
}

func (r *ServiceRepository) Service(orgID, serviceName string) (*domain.Service, error) {
	result := Service{}
	if err := r.DB.Where(&Service{OrgID: orgID, Name: serviceName}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from serviecs: %v", err)
	}

	service := result.Domain()
	return &service, nil
}

func (r *ServiceRepository) Exists(orgID, serviceName string) bool {
	var count int64
	if err := r.DB.Model(&Service{}).Where(&Service{OrgID: orgID, Name: serviceName}).Count(&count).Error; err != nil {
		return false // FIXME Add error message
	}

	if count == 0 {
		return false
	}

	return true
}

func (r *ServiceRepository) Delete(orgID, serviceName string) error {
	if err := r.DB.Delete(&Service{OrgID: orgID, Name: serviceName}).Error; err != nil {
		return fmt.Errorf("delete from services: %v", err)
	}

	return nil
}
