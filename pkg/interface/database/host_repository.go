package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
)

type HostRepository struct {
	DB *gorm.DB
}

type Host struct {
	OrgID            string `gorm:"column:org_id;            type:varchar(16);  not null;"`
	ID               string `gorm:"column:id;                type:varchar(16);  not null; primary key"`
	Name             string `gorm:"column:name;              type:varchar(128);  not null;"`
	Status           string `gorm:"column:status;            type:enum('working', 'standby', 'maintenance', 'poweroff');  not null;"`
	Memo             string `gorm:"column:memo;              type:varchar(128);  not null; default:''"`
	DisplayName      string `gorm:"column:display_name;      type:varchar(128);"`
	CustomIdentifier string `gorm:"column:custom_identifier; type:varchar(128);"`
	CreatedAt        int64  `gorm:"column:created_at;        type:bigint;"`
	RetiredAt        int64  `gorm:"column:retired_at;        type:bigint;"`
	IsRetired        bool   `gorm:"column:is_retired;        type:boolean;"`
	Roles            string `gorm:"column:roles;             type:text;"`
	RoleFullNames    string `gorm:"column:role_fullnames;    type:text;"`
	Interfaces       string `gorm:"column:interfaces;        type:text;"`
	Checks           string `gorm:"column:checks;            type:text;"`
	Meta             string `gorm:"column:meta;              type:text;"`
}

func (h Host) Domain() (domain.Host, error) {
	host := domain.Host{
		OrgID:            h.OrgID,
		ID:               h.ID,
		Name:             h.Name,
		Status:           h.Status,
		Memo:             h.Memo,
		DisplayName:      h.DisplayName,
		CustomIdentifier: h.CustomIdentifier,
		CreatedAt:        h.CreatedAt,
		RetiredAt:        h.RetiredAt,
		IsRetired:        h.IsRetired,
	}

	if err := json.Unmarshal([]byte(h.Roles), &host.Roles); err != nil {
		return host, fmt.Errorf("unmarshal host.Roles: %v", err)
	}

	if err := json.Unmarshal([]byte(h.RoleFullNames), &host.RoleFullNames); err != nil {
		return host, fmt.Errorf("unmarshal host.RoleFullNames: %v", err)
	}

	if err := json.Unmarshal([]byte(h.Interfaces), &host.Interfaces); err != nil {
		return host, fmt.Errorf("unmarshal host.Interfaces: %v", err)
	}

	if err := json.Unmarshal([]byte(h.Checks), &host.Checks); err != nil {
		return host, fmt.Errorf("unmarshal host.Checks: %v", err)
	}

	if err := json.Unmarshal([]byte(h.Meta), &host.Meta); err != nil {
		return host, fmt.Errorf("unmarshal host.Meta: %v", err)
	}

	return host, nil
}

func NewHostRepository(handler SQLHandler) *HostRepository {
	db, err := gorm.Open(handler.Dialect(), handler.Raw())
	if err != nil {
		panic(err)
	}
	db.LogMode(handler.IsDebugging())

	if err := db.AutoMigrate(&Host{}).Error; err != nil {
		panic(fmt.Errorf("auto migrate host: %v", err))
	}

	return &HostRepository{
		DB: db,
	}
}

func (repo *HostRepository) ActiveList(orgID string) (*domain.Hosts, error) {
	result := make([]Host, 0)
	if err := repo.DB.Where(&Host{OrgID: orgID, IsRetired: false}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from hosts: %v", err)
	}

	out := make([]domain.Host, 0)
	for _, r := range result {
		host, err := r.Domain()
		if err != nil {
			return nil, fmt.Errorf("domain: %v", err)
		}

		out = append(out, host)
	}

	return &domain.Hosts{Hosts: out}, nil
}

func (repo *HostRepository) List(orgID string) (*domain.Hosts, error) {
	result := make([]Host, 0)
	if err := repo.DB.Where(&Host{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from hosts: %v", err)
	}

	out := make([]domain.Host, 0)
	for _, r := range result {
		host, err := r.Domain()
		if err != nil {
			return nil, fmt.Errorf("domain: %v", err)
		}

		out = append(out, host)
	}

	return &domain.Hosts{Hosts: out}, nil
}

func (repo *HostRepository) Save(orgID string, host *domain.Host) (*domain.HostID, error) {
	roles, err := json.Marshal(host.Roles)
	if err != nil {
		return nil, fmt.Errorf("marshal host.Roles: %v", err)
	}

	roleFullnames, err := json.Marshal(host.RoleFullNames)
	if err != nil {
		return nil, fmt.Errorf("marshal host.RoleFullNames: %v", err)
	}

	interfaces, err := json.Marshal(host.Interfaces)
	if err != nil {
		return nil, fmt.Errorf("marshal host.Interfaces: %v", err)
	}

	checks, err := json.Marshal(host.Checks)
	if err != nil {
		return nil, fmt.Errorf("marshal host.Checks: %v", err)
	}

	meta, err := json.Marshal(host.Meta)
	if err != nil {
		return nil, fmt.Errorf("marshal host.Meta: %v", err)
	}

	update := Host{
		Name:             host.Name,
		Status:           host.Status,
		Memo:             host.Memo,
		DisplayName:      host.DisplayName,
		CustomIdentifier: host.CustomIdentifier,
		CreatedAt:        host.CreatedAt,
		RetiredAt:        host.RetiredAt,
		IsRetired:        host.IsRetired,
		Roles:            string(roles),
		RoleFullNames:    string(roleFullnames),
		Interfaces:       string(interfaces),
		Checks:           string(checks),
		Meta:             string(meta),
	}

	if err := repo.DB.Where(&Host{OrgID: orgID, ID: host.ID}).Assign(&update).FirstOrCreate(&Host{}).Error; err != nil {
		return nil, fmt.Errorf("first or create: %v", err)
	}

	return &domain.HostID{ID: host.ID}, nil
}

func (repo *HostRepository) Host(orgID, hostID string) (*domain.Host, error) {
	result := Host{}
	if err := repo.DB.Where(&Host{OrgID: orgID, ID: hostID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from hosts: %v", err)
	}

	host, err := result.Domain()
	if err != nil {
		return nil, fmt.Errorf("domain: %v", err)
	}

	return &host, nil
}

func (repo *HostRepository) Exists(orgID, hostID string) bool {
	if repo.DB.Where(&Host{OrgID: orgID, ID: hostID}).First(&Host{}).RecordNotFound() {
		return false
	}

	return true
}

func (repo *HostRepository) Status(orgID, hostID, status string) (*domain.Success, error) {
	if err := repo.DB.Where(&Host{OrgID: orgID, ID: hostID}).Assign(&Host{Status: status}).FirstOrCreate(&Host{}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("update hosts: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) Retire(orgID, hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	if err := repo.DB.Where(&Host{OrgID: orgID, ID: hostID}).Assign(&Host{IsRetired: true, RetiredAt: time.Now().Unix()}).FirstOrCreate(&Host{}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("update hosts: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) SaveRoleFullNames(orgID, hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	roles, err := json.Marshal(names.Roles())
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	roleFullnames, err := json.Marshal(names.Names)
	if err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("marshal: %v", err)
	}

	if err := repo.DB.Where(&Host{OrgID: orgID, ID: hostID}).Assign(&Host{RoleFullNames: string(roleFullnames), Roles: string(roles)}).FirstOrCreate(&Host{}).Error; err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}
