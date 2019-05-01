package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostRepository interface {
	SaveCustomGraphDefs(defs domain.CustomGraphDefs) error
	MetricNames() ([]string, error)
	MetricValuesLatest(hostID, metricName []string) (domain.HostMetricValues, error)
	MetricValues(hostID, metricName string, from, to int64) (domain.HostMetricValues, error)
	SaveMetricValues(v domain.HostMetricValues) error
	Exists(hostName string) bool
	FindByID(hostID string) (domain.Host, error)
	FindByName(hostName string) (domain.Host, error)
	FindAll() (domain.Hosts, error)
	Save(host domain.Host) error
	DeleteByID(hostID string) error
	DeleteByName(hostName string) error
}
