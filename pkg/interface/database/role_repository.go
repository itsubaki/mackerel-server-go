package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type RoleRepository struct {
	DB *gorm.DB
}

type Role struct {
	OrgID       string `gorm:"column:org_id; type:varchar(16);        not null; primary_key"`
	ServiceName string `gorm:"column:service_name; type:varchar(128); not null; primary_key"`
	Name        string `gorm:"column:name;   type:varchar(128);       not null; primary_key"`
	Memo        string `gorm:"column:memo;   type:varchar(128);       not null; default:''"`
}

func NewRoleRepository(handler SQLHandler) *RoleRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

	if err := db.AutoMigrate(&Role{}).AddForeignKey("org_id, service_name", "services(org_id, name)", "CASCADE", "CASCADE").Error; err != nil {
		panic(fmt.Errorf("auto migrate role: %v", err))
	}

	return &RoleRepository{
		DB: db,
	}
}

func (repo *RoleRepository) List(orgID string) (map[string][]string, error) {
	result := make([]Role, 0)
	if err := repo.DB.Where(&Role{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from roles: %v", err)
	}

	roles := make(map[string][]string)
	for _, r := range result {
		if _, ok := roles[r.ServiceName]; !ok {
			roles[r.ServiceName] = make([]string, 0)
		}

		roles[r.ServiceName] = append(roles[r.ServiceName], r.Name)
	}

	return roles, nil
}

func (repo *RoleRepository) ListWith(orgID, serviceName string) (*domain.Roles, error) {
	result := make([]Role, 0)
	if err := repo.DB.Where(&Role{OrgID: orgID, ServiceName: serviceName}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from roles: %v", err)
	}

	out := make([]domain.Role, 0)
	for _, r := range result {
		out = append(out, domain.Role{
			OrgID:       r.OrgID,
			ServiceName: r.ServiceName,
			Name:        r.Name,
			Memo:        r.Memo,
		})
	}

	return &domain.Roles{Roles: out}, nil
}

func (repo *RoleRepository) Role(orgID, serviceName, roleName string) (*domain.Role, error) {
	result := Role{}
	if err := repo.DB.Where(&Role{OrgID: orgID, ServiceName: serviceName, Name: roleName}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from roles: %v", err)
	}

	role := domain.Role{
		OrgID:       result.OrgID,
		ServiceName: result.ServiceName,
		Name:        result.Name,
		Memo:        result.Memo,
	}

	return &role, nil
}

func (repo *RoleRepository) Save(orgID, serviceName string, r *domain.Role) error {
	update := Role{
		OrgID:       orgID,
		ServiceName: serviceName,
		Name:        r.Name,
		Memo:        r.Memo,
	}

	if err := repo.DB.Where(&Role{OrgID: orgID, ServiceName: serviceName}).Assign(&update).FirstOrCreate(&Role{}).Error; err != nil {
		return fmt.Errorf("insert into roles: %v", err)
	}

	return nil
}

func (repo *RoleRepository) Exists(orgID, serviceName, roleName string) bool {
	if repo.DB.Where(&Role{OrgID: orgID, ServiceName: serviceName, Name: roleName}).First(&Role{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *RoleRepository) Delete(orgID, serviceName, roleName string) error {
	if err := repo.DB.Delete(&Role{OrgID: orgID, ServiceName: serviceName, Name: roleName}).Error; err != nil {
		return fmt.Errorf("delete from roles: %v", err)
	}

	return nil
}
