package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostRepository interface {
	List(orgID string) (*domain.Hosts, error)
	Save(orgID string, host *domain.Host) (*domain.HostID, error)
	Host(orgID, hostID string) (*domain.Host, error)
	Exists(orgID, hostID string) bool
	Status(orgID, hostID, status string) (*domain.Success, error)
	SaveRoleFullNames(orgID, hostID string, names *domain.RoleFullNames) (*domain.Success, error)
	Retire(orgID, hostID string, retire *domain.HostRetire) (*domain.Success, error)

	ExistsMetric(orgID, hostID, name string) bool
	MetricNames(orgID, hostID string) (*domain.MetricNames, error)
	MetricValues(orgID, hostID, name string, from, to int64) (*domain.MetricValues, error)
	MetricValuesLimit(orgID, hostID, name string, limit int) (*domain.MetricValues, error)
	MetricValuesLatest(orgID string, hostID, name []string) (*domain.TSDBLatest, error)
	SaveMetricValues(orgID string, values []domain.MetricValue) (*domain.Success, error)

	ExistsMetadata(orgID, hostID, namespace string) bool
	MetadataList(orgID, hostID string) (*domain.HostMetadataList, error)
	Metadata(orgID, hostID, namespace string) (interface{}, error)
	SaveMetadata(orgID, hostID, namespace string, metadata interface{}) (*domain.Success, error)
	DeleteMetadata(orgID, hostID, namespace string) (*domain.Success, error)
}
