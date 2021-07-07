package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceMetaRepository struct {
	DB *gorm.DB
}

type ServiceMeta struct {
	OrgID       string `gorm:"column:org_id;       type:varchar(16);  not null; primary_key"`
	ServiceName string `gorm:"column:service_name; type:varchar(16);  not null; primary_key"`
	Namespace   string `gorm:"column:namespace;    type:varchar(128); not null; primary_key"`
	Metadata    string `gorm:"column:meta;         type:text"`
}

func NewServiceMetaRepository(handler SQLHandler) *ServiceMetaRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&ServiceMeta{}); err != nil {
		panic(fmt.Errorf("auto migrate service_meta: %v", err))
	}

	return &ServiceMetaRepository{
		DB: db,
	}
}

func (r *ServiceMetaRepository) Exists(orgID, serviceName, namespace string) bool {
	if err := r.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).First(&ServiceMeta{}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (r *ServiceMetaRepository) List(orgID, serviceName string) (*domain.ServiceMetadataList, error) {
	result := make([]ServiceMeta, 0)
	if err := r.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from service_meta: %#v", err)
	}

	out := make([]domain.ServiceMetadata, 0)
	for _, r := range result {
		out = append(out, domain.ServiceMetadata{
			Namespace: r.Namespace,
		})
	}

	return &domain.ServiceMetadataList{Metadata: out}, nil
}

func (r *ServiceMetaRepository) Metadata(orgID, serviceName, namespace string) (interface{}, error) {
	result := ServiceMeta{}
	if err := r.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from serviec_meta: %v", err)
	}

	var out interface{}
	if err := json.Unmarshal([]byte(result.Metadata), &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return out, nil
}

func (r *ServiceMetaRepository) Save(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := r.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).Assign(&ServiceMeta{Metadata: string(meta)}).FirstOrCreate(&ServiceMeta{}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("first or create: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (r *ServiceMetaRepository) Delete(orgID, serviceName, namespace string) (*domain.Success, error) {
	if err := r.DB.Delete(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("delete from service_meta: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
