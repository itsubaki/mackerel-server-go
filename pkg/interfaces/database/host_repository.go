package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	SQLHanler SQLHandler
}

// select * from hosts
func (repo *HostRepository) List() (*domain.Hosts, error) {
	return nil, nil
}

// insert into hosts values(${host.ID}, ...)
func (repo *HostRepository) Save(host *domain.Host) (*domain.HostID, error) {
	return &domain.HostID{}, nil
}

// select * from hosts where ID=${hostID}
func (repo *HostRepository) Host(hostID string) (*domain.Host, error) {
	return nil, fmt.Errorf("host not found")
}

// select * from hosts where id=${hostID} limit=1
func (repo *HostRepository) Exists(hostID string) bool {
	return false
}

// update hosts set status=${status} where id=${hostID}
func (repo *HostRepository) Status(hostID, status string) (*domain.Success, error) {
	return &domain.Success{Success: false}, fmt.Errorf("host not found")
}

// update hosts set roles=${roles} where id=${hostID}
func (repo *HostRepository) SaveRoleFullNames(hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	return nil, fmt.Errorf("host not found")
}

// update hosts set isRetired=true, retiredAt=time.Now().Unix() where id=${hostID}
func (repo *HostRepository) Retire(hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) ExistsMetric(hostID, name string) bool {
	return false
}

func (repo *HostRepository) MetricNames(hostID string) (*domain.MetricNames, error) {
	return &domain.MetricNames{}, nil
}

func (repo *HostRepository) MetricValues(hostID, name string, from, to int) (*domain.MetricValues, error) {
	return &domain.MetricValues{}, nil
}

func (repo *HostRepository) MetricValuesLatest(hostId, name []string) (*domain.TSDBLatest, error) {
	return &domain.TSDBLatest{}, nil
}

func (repo *HostRepository) SaveMetricValues(values []domain.MetricValue) (*domain.Success, error) {
	return nil, nil
}

func (repo *HostRepository) ExistsMetadata(hostID, namespace string) bool {
	return true
}

func (repo *HostRepository) MetadataList(hostID string) (*domain.HostMetadataList, error) {
	return &domain.HostMetadataList{}, nil
}

func (repo *HostRepository) Metadata(hostID, namespace string) (interface{}, error) {
	return "", nil
}

func (repo *HostRepository) SaveMetadata(hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) DeleteMetadata(hostID, namespace string) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}
