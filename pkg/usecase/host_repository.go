package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostRepository interface {
	List(org string) (*domain.Hosts, error)
	Save(org string, host *domain.Host) (*domain.HostID, error)
	Host(org, hostID string) (*domain.Host, error)
	Exists(org, hostID string) bool
	Status(org, hostID, status string) (*domain.Success, error)
	SaveRoleFullNames(org, hostID string, names *domain.RoleFullNames) (*domain.Success, error)
	Retire(org, hostID string, retire *domain.HostRetire) (*domain.Success, error)

	ExistsMetric(org, hostID, name string) bool
	MetricNames(org, hostID string) (*domain.MetricNames, error)
	MetricValues(org, hostID, name string, from, to int64) (*domain.MetricValues, error)
	MetricValuesLimit(org, hostID, name string, limit int) (*domain.MetricValues, error)
	MetricValuesLatest(org string, hostID, name []string) (*domain.TSDBLatest, error)
	SaveMetricValues(org string, values []domain.MetricValue) (*domain.Success, error)

	ExistsMetadata(org, hostID, namespace string) bool
	MetadataList(org, hostID string) (*domain.HostMetadataList, error)
	Metadata(org, hostID, namespace string) (interface{}, error)
	SaveMetadata(org, hostID, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteMetadata(org, hostID, namespace string) (*domain.Success, error)
}
