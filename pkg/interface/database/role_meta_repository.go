package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RoleMetaRepository struct {
	DB *gorm.DB
}

type RoleMeta struct {
	OrgID       string `gorm:"column:org_id;       type:varchar(16);  not null; primary_key"`
	ServiceName string `gorm:"column:service_name; type:varchar(16);  not null"`
	RoleName    string `gorm:"column:role_name;    type:varchar(16);  not null; primary_key"`
	Namespace   string `gorm:"column:namespace;    type:varchar(128); not null; primary_key"`
	Metadata    string `gorm:"column:meta;         type:text"`
}

func NewRoleMetaRepository(handler SQLHandler) *RoleMetaRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&RoleMeta{}); err != nil {
		panic(fmt.Errorf("auto migrate role_meta: %v", err))
	}

	return &RoleMetaRepository{
		DB: db,
	}
}

func (r *RoleMetaRepository) Exists(orgID, serviceName, roleName, namespace string) bool {
	if err := r.DB.Where(&RoleMeta{OrgID: orgID, ServiceName: serviceName, RoleName: roleName, Namespace: namespace}).First(&RoleMeta{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (r *RoleMetaRepository) List(orgID, serviceName, roleName string) (*domain.RoleMetadataList, error) {
	result := make([]RoleMeta, 0)
	if err := r.DB.Where(&RoleMeta{OrgID: orgID, ServiceName: serviceName, RoleName: roleName}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from role_meta: %v", err)
	}

	out := make([]domain.RoleMetadata, 0)
	for _, r := range result {
		out = append(out, domain.RoleMetadata{
			Namespace: r.Namespace,
		})
	}

	return &domain.RoleMetadataList{Metadata: out}, nil
}

func (r *RoleMetaRepository) Metadata(orgID, serviceName, roleName, namespace string) (interface{}, error) {
	result := RoleMeta{}
	if err := r.DB.Where(&RoleMeta{OrgID: orgID, ServiceName: serviceName, RoleName: roleName, Namespace: namespace}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from role_meta: %v", err)
	}

	var out interface{}
	if err := json.Unmarshal([]byte(result.Metadata), &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return out, nil
}

func (r *RoleMetaRepository) Save(orgID, serviceName, roleName, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := r.DB.Where(&RoleMeta{OrgID: orgID, ServiceName: serviceName, RoleName: roleName, Namespace: namespace}).Assign(&RoleMeta{Metadata: string(meta)}).FirstOrCreate(&RoleMeta{}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("first or create: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (r *RoleMetaRepository) Delete(orgID, serviceName, roleName, namespace string) (*domain.Success, error) {
	if err := r.DB.Delete(&RoleMeta{OrgID: orgID, ServiceName: serviceName, RoleName: roleName, Namespace: namespace}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("delete from role_meta: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
