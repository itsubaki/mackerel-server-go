package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	DB SQLHandler
}

func NewHostRepository(handler SQLHandler) *HostRepository {
	return &HostRepository{
		DB: handler,
	}
}

// select * from hosts
func (repo *HostRepository) List() (*domain.Hosts, error) {
	return nil, nil
}

// insert into hosts values(${name}, ${meta}, ${interfaces}, ${checks}, ${display_name}, ${custom_identifier}, ${created_at}, ${id}, ${status}, ${memo}, ${roles}, ${is_retired}, ${retired_at} )
func (repo *HostRepository) Save(host *domain.Host) (*domain.HostID, error) {
	return &domain.HostID{}, nil
}

// select * from hosts where id=${hostID}
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

// update hosts set is_retired=true, retired_at=time.Now().Unix() where id=${hostID}
func (repo *HostRepository) Retire(hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}

// select * from host_metrics where host_id=${hostID} and name=${name} limit=1
func (repo *HostRepository) ExistsMetric(hostID, name string) bool {
	return false
}

// select distinct name from host_metrics where host_id=${hostID}
func (repo *HostRepository) MetricNames(hostID string) (*domain.MetricNames, error) {
	return &domain.MetricNames{}, nil
}

// select value from host_metric_values where host_id=${hostID} and name=${name} and ${from} < from and to < ${to}
func (repo *HostRepository) MetricValues(hostID, name string, from, to int) (*domain.MetricValues, error) {
	return &domain.MetricValues{}, nil
}

// select * from host_metric_values_latest where host_id=${hostID} and name=${name}
func (repo *HostRepository) MetricValuesLatest(hostID, name []string) (*domain.TSDBLatest, error) {
	return &domain.TSDBLatest{}, nil
}

// insert into host_metric_values values(${host_id}, ${name}, ${time}, ${value})
func (repo *HostRepository) SaveMetricValues(values []domain.MetricValue) (*domain.Success, error) {
	return nil, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace} limit=1
func (repo *HostRepository) ExistsMetadata(hostID, namespace string) bool {
	return true
}

// select namespace from host_metadata where host_id=${hostID}
func (repo *HostRepository) MetadataList(hostID string) (*domain.HostMetadataList, error) {
	return &domain.HostMetadataList{}, nil
}

// select * from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostRepository) Metadata(hostID, namespace string) (interface{}, error) {
	return "", nil
}

// insert into host_metadata values(${hostID}, ${namespace}, ${metadata})
func (repo *HostRepository) SaveMetadata(hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}

// delete from host_metadata where host_id=${hostID} and namespace=${namespace}
func (repo *HostRepository) DeleteMetadata(hostID, namespace string) (*domain.Success, error) {
	return &domain.Success{Success: true}, nil
}
