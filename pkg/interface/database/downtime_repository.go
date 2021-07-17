package database

import (
	"encoding/json"
	"fmt"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ usecase.DowntimeRepository = (*DowntimeRepository)(nil)

type DowntimeRepository struct {
	DB *gorm.DB
}

type Downtime struct {
	OrgID                string `gorm:"column:org_id;                 type:varchar(16);  not null"`
	ID                   string `gorm:"column:id;                     type:varchar(128); not null; primary_key"`
	Name                 string `gorm:"column:name;                   type:varchar(128); not null"`
	Memo                 string `gorm:"column:memo;                   type:text;"`
	Start                int64  `gorm:"column:start;                  type:bigint;"`
	Duration             int64  `gorm:"column:duration;               type:bigint;"`
	Recurrence           string `gorm:"column:recurrence;             type:text;"`
	ServiceScopes        string `gorm:"column:service_scopes;         type:text;"`
	ServiceExcludeScopes string `gorm:"column:service_exclude_scopes; type:text;"`
	RoleScopes           string `gorm:"column:role_scopes;            type:text;"`
	RoleExcludeScopes    string `gorm:"column:role_exclude_scopes;    type:text;"`
	MonitorScopes        string `gorm:"column:monitor_scopes;         type:text;"`
	MonitorExcludeScopes string `gorm:"column:monitor_exclude_scopes; type:text;"`
}

func (t Downtime) Domain() (domain.Downtime, error) {
	downtime := domain.Downtime{
		OrgID:    t.OrgID,
		ID:       t.ID,
		Name:     t.Name,
		Memo:     t.Memo,
		Start:    t.Start,
		Duration: t.Duration,
	}

	if err := json.Unmarshal([]byte(t.Recurrence), &downtime.Recurrence); err != nil {
		return downtime, fmt.Errorf("unmarshal downitme.Recurrence: %v", err)
	}

	if err := json.Unmarshal([]byte(t.ServiceScopes), &downtime.ServiceScopes); err != nil {
		return downtime, fmt.Errorf("unmarshal downitme.ServiceScopes: %v", err)
	}

	if err := json.Unmarshal([]byte(t.ServiceExcludeScopes), &downtime.ServiceExcludeScopes); err != nil {
		return downtime, fmt.Errorf("unmarshal downitme.ServiceExcludeScopes: %v", err)
	}

	if err := json.Unmarshal([]byte(t.RoleScopes), &downtime.RoleScopes); err != nil {
		return downtime, fmt.Errorf("unmarshal downitme.RoleScopes: %v", err)
	}

	if err := json.Unmarshal([]byte(t.RoleExcludeScopes), &downtime.RoleExcludeScopes); err != nil {
		return downtime, fmt.Errorf("unmarshal downitme.RoleExcludeScopes: %v", err)
	}

	if err := json.Unmarshal([]byte(t.MonitorScopes), &downtime.MonitorScopes); err != nil {
		return downtime, fmt.Errorf("unmarshal downitme.MonitorScopes: %v", err)
	}

	if err := json.Unmarshal([]byte(t.MonitorExcludeScopes), &downtime.MonitorExcludeScopes); err != nil {
		return downtime, fmt.Errorf("unmarshal downitme.MonitorExcludeScopes: %v", err)
	}

	return downtime, nil
}

func NewDowntimeRepository(handler SQLHandler) *DowntimeRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&Downtime{}); err != nil {
		panic(fmt.Errorf("auto migrate downtime: %v", err))
	}

	return &DowntimeRepository{
		DB: db,
	}
}

func (r *DowntimeRepository) List(orgID string) (*domain.Downtimes, error) {
	result := make([]Downtime, 0)
	if err := r.DB.Where(&Downtime{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from downtimes: %v", err)
	}

	out := make([]domain.Downtime, 0)
	for _, r := range result {
		t, err := r.Domain()
		if err != nil {
			return nil, fmt.Errorf("domain: %v", err)
		}

		out = append(out, t)
	}

	return &domain.Downtimes{Downtimes: out}, nil
}

func (r *DowntimeRepository) Save(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	recurrence, err := json.Marshal(downtime.Recurrence)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.Recurrenc: %v", err)
	}

	serviceScopes, err := json.Marshal(downtime.ServiceScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.ServiceScopes: %v", err)
	}

	serviceExcludeScopes, err := json.Marshal(downtime.ServiceExcludeScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.ServiceExcludeScopes: %v", err)
	}

	roleScopes, err := json.Marshal(downtime.RoleScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.RoleScopes: %v", err)
	}

	roleExcludeScopes, err := json.Marshal(downtime.RoleExcludeScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.RoleExcludeScopes: %v", err)
	}

	monitorScopes, err := json.Marshal(downtime.MonitorScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.MonitorScopes: %v", err)
	}

	monitorExcludeScopes, err := json.Marshal(downtime.MonitorExcludeScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.MonitorExcludeScopes: %v", err)
	}

	create := Downtime{
		OrgID:                orgID,
		ID:                   downtime.ID,
		Name:                 downtime.Name,
		Memo:                 downtime.Memo,
		Start:                downtime.Start,
		Duration:             downtime.Duration,
		Recurrence:           string(recurrence),
		ServiceScopes:        string(serviceScopes),
		ServiceExcludeScopes: string(serviceExcludeScopes),
		RoleScopes:           string(roleScopes),
		RoleExcludeScopes:    string(roleExcludeScopes),
		MonitorScopes:        string(monitorScopes),
		MonitorExcludeScopes: string(monitorExcludeScopes),
	}

	if err := r.DB.Create(&create).Error; err != nil {
		return nil, fmt.Errorf("insert into downtimes: %v", err)
	}

	return downtime, nil
}

func (r *DowntimeRepository) Update(orgID string, downtime *domain.Downtime) (*domain.Downtime, error) {
	recurrence, err := json.Marshal(downtime.Recurrence)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.Recurrence: %v", err)
	}

	serviceScopes, err := json.Marshal(downtime.ServiceScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.ServiceScopes: %v", err)
	}

	serviceExcludeScopes, err := json.Marshal(downtime.ServiceExcludeScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.ServiceExcludeScopes: %v", err)
	}

	roleScopes, err := json.Marshal(downtime.RoleScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.RoleScopes: %v", err)
	}

	roleExcludeScopes, err := json.Marshal(downtime.RoleExcludeScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.RoleExcludeScopes: %v", err)
	}

	monitorScopes, err := json.Marshal(downtime.MonitorScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.MonitorScopes: %v", err)
	}

	monitorExcludeScopes, err := json.Marshal(downtime.MonitorExcludeScopes)
	if err != nil {
		return nil, fmt.Errorf("marshal downtime.MonitorExcludeScopes: %v", err)
	}

	update := Downtime{
		Name:                 downtime.Name,
		Memo:                 downtime.Memo,
		Start:                downtime.Start,
		Duration:             downtime.Duration,
		Recurrence:           string(recurrence),
		ServiceScopes:        string(serviceScopes),
		ServiceExcludeScopes: string(serviceExcludeScopes),
		RoleScopes:           string(roleScopes),
		RoleExcludeScopes:    string(roleExcludeScopes),
		MonitorScopes:        string(monitorScopes),
		MonitorExcludeScopes: string(monitorExcludeScopes),
	}

	if err := r.DB.Where(&Downtime{OrgID: orgID, ID: downtime.ID}).Assign(&update).FirstOrCreate(&Downtime{}).Error; err != nil {
		return nil, fmt.Errorf("first or create: %v", err)
	}

	return downtime, nil
}

func (r *DowntimeRepository) Downtime(orgID, downtimeID string) (*domain.Downtime, error) {
	result := Downtime{}
	if err := r.DB.Where(&Downtime{OrgID: orgID, ID: downtimeID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from downitmes: %v", err)
	}

	downtime, err := result.Domain()
	if err != nil {
		return nil, fmt.Errorf("domain: %v", err)
	}

	return &downtime, nil
}

func (r *DowntimeRepository) Delete(orgID, downtimeID string) (*domain.Downtime, error) {
	result := Downtime{}
	if err := r.DB.Where(&Downtime{OrgID: orgID, ID: downtimeID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from downitmes: %v", err)
	}

	if err := r.DB.Delete(&Downtime{OrgID: orgID, ID: downtimeID}).Error; err != nil {
		return nil, fmt.Errorf("delete from downtimes: %v", err)
	}

	downtime, err := result.Domain()
	if err != nil {
		return nil, fmt.Errorf("domain: %v", err)
	}

	return &downtime, nil
}
