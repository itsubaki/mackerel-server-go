package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
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
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

	if err := db.AutoMigrate(&ServiceMeta{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate service_meta: %v", err))
	}

	return &ServiceMetaRepository{
		DB: db,
	}
}

func (repo *ServiceMetaRepository) Exists(orgID, serviceName, namespace string) bool {
	if repo.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).First(&ServiceMeta{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *ServiceMetaRepository) List(orgID, serviceName string) (*domain.ServiceMetadataList, error) {
	result := make([]ServiceMeta, 0)
	if err := repo.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName}).Find(&result).Error; err != nil {
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

func (repo *ServiceMetaRepository) Metadata(orgID, serviceName, namespace string) (interface{}, error) {
	result := ServiceMeta{}
	if err := repo.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from serviec_meta: %v", err)
	}

	var out interface{}
	if err := json.Unmarshal([]byte(result.Metadata), &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return out, nil
}

func (repo *ServiceMetaRepository) Save(orgID, serviceName, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.DB.Where(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).Assign(&ServiceMeta{Metadata: string(meta)}).FirstOrCreate(&ServiceMeta{}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("first or create: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (repo *ServiceMetaRepository) Delete(orgID, serviceName, namespace string) (*domain.Success, error) {
	if err := repo.DB.Delete(&ServiceMeta{OrgID: orgID, ServiceName: serviceName, Namespace: namespace}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("delete from service_meta: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
