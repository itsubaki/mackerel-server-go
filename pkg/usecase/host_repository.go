package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostRepository interface {
	List() (*domain.Hosts, error)
	Save(host *domain.Host) (*domain.HostID, error)
	Host(hostID string) (*domain.Host, error)
	Exists(hostID string) bool
	Status(hostID, status string) (*domain.Success, error)
	SaveRoleFullNames(hostID string, names *domain.RoleFullNames) (*domain.Success, error)
	Retire(hostID string, retire *domain.HostRetire) (*domain.Success, error)

	ExistsMetric(hostID, name string) bool
	MetricNames(hostID string) (*domain.MetricNames, error)
	MetricValues(hostID, name string, from, to int) (*domain.MetricValues, error)
	MetricValuesLatest(hostID, name []string) (*domain.TSDBLatest, error)
	SaveMetricValues(values []domain.MetricValue) (*domain.Success, error)

	ExistsMetadata(hostID, namespace string) bool
	MetadataList(hostID string) (*domain.HostMetadataList, error)
	Metadata(hostID, namespace string) (interface{}, error)
	SaveMetadata(hostID, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteMetadata(hostID, namespace string) (*domain.Success, error)
}
