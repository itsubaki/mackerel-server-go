package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type HostMetaRepository struct {
	DB *gorm.DB
}

type HostMeta struct {
	OrgID     string `gorm:"column:org_id;    type:varchar(16);  not null;"`
	HostID    string `gorm:"column:host_id;   type:varchar(16);  not null; primary_key"`
	Namespace string `gorm:"column:namespace; type:varchar(128); not null; primary_key"`
	Metadata  string `gorm:"column:meta;      type:text"`
}

func NewHostMetaRepository(handler SQLHandler) *HostMetaRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())
	if err := db.AutoMigrate(&HostMeta{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate host_meta: %v", err))
	}

	return &HostMetaRepository{
		DB: db,
	}
}

func (repo *HostMetaRepository) Exists(orgID, hostID, namespace string) bool {
	if repo.DB.Where(&HostMeta{OrgID: orgID, HostID: hostID, Namespace: namespace}).First(&HostMeta{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *HostMetaRepository) List(orgID, hostID string) (*domain.HostMetadataList, error) {
	result := make([]HostMeta, 0)
	if err := repo.DB.Where(&HostMeta{OrgID: orgID, HostID: hostID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from host_meta: %v", err)
	}

	out := make([]domain.Namespace, 0)
	for _, r := range result {
		out = append(out, domain.Namespace{
			Namespace: r.Namespace,
		})
	}

	return &domain.HostMetadataList{Metadata: out}, nil
}

func (repo *HostMetaRepository) Metadata(orgID, hostID, namespace string) (interface{}, error) {
	result := HostMeta{}
	if err := repo.DB.Where(&HostMeta{OrgID: orgID, HostID: hostID, Namespace: namespace}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from host_meta: %v", err)
	}

	var out interface{}
	if err := json.Unmarshal([]byte(result.Metadata), &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return out, nil
}

func (repo *HostMetaRepository) Save(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	meta, err := json.Marshal(metadata)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.DB.Where(&HostMeta{OrgID: orgID, HostID: hostID, Namespace: namespace}).Assign(&HostMeta{Metadata: string(meta)}).FirstOrCreate(&HostMeta{}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("first or create: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (repo *HostMetaRepository) Delete(orgID, hostID, namespace string) (*domain.Success, error) {
	if err := repo.DB.Delete(&HostMeta{OrgID: orgID, HostID: hostID, Namespace: namespace}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("delete from host_meta: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
