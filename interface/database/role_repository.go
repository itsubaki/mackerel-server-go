package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.RoleRepository = (*RoleRepository)(nil)

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
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&Role{}); err != nil {
		panic(fmt.Errorf("auto migrate role: %v", err))
	}

	return &RoleRepository{
		DB: db,
	}
}

func (r *RoleRepository) List(orgID string) (map[string][]string, error) {
	result := make([]Role, 0)
	if err := r.DB.Where(&Role{OrgID: orgID}).Find(&result).Error; err != nil {
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

func (r *RoleRepository) ListWith(orgID, serviceName string) (*domain.Roles, error) {
	result := make([]Role, 0)
	if err := r.DB.Where(&Role{OrgID: orgID, ServiceName: serviceName}).Find(&result).Error; err != nil {
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

func (r *RoleRepository) Role(orgID, serviceName, roleName string) (*domain.Role, error) {
	result := Role{}
	if err := r.DB.Where(&Role{OrgID: orgID, ServiceName: serviceName, Name: roleName}).First(&result).Error; err != nil {
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

func (r *RoleRepository) Save(orgID, serviceName string, role *domain.Role) error {
	update := Role{
		OrgID:       orgID,
		ServiceName: serviceName,
		Name:        role.Name,
		Memo:        role.Memo,
	}

	if err := r.DB.Where(&Role{OrgID: orgID, ServiceName: serviceName, Name: role.Name}).Assign(&update).FirstOrCreate(&Role{}).Error; err != nil {
		return fmt.Errorf("insert into roles: %v", err)
	}

	return nil
}

func (r *RoleRepository) Exists(orgID, serviceName, roleName string) bool {
	var count int64
	if err := r.DB.Model(&Role{}).Where(&Role{OrgID: orgID, ServiceName: serviceName, Name: roleName}).Count(&count).Error; err != nil {
		return false // FIXME Add error message
	}

	if count == 0 {
		return false
	}

	return true
}

func (r *RoleRepository) Delete(orgID, serviceName, roleName string) error {
	if err := r.DB.Delete(&Role{OrgID: orgID, ServiceName: serviceName, Name: roleName}).Error; err != nil {
		return fmt.Errorf("delete from roles: %v", err)
	}

	return nil
}
