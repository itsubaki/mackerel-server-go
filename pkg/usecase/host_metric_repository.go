package usecase

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetricRepository interface {
	FindAll() (domain.HostMetricValues, error)
	Latest(hostID, metricName []string) (domain.HostMetricValues, error)
	ExistsByName(hostID, metricName string) bool
	FindByID(hostID string) (domain.HostMetricValues, error)
	FindBy(hostID, metricName string, from, to int64) (domain.HostMetricValues, error)
	Save(v domain.HostMetricValues) error
}
