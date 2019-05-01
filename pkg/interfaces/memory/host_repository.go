package memory

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	Hosts            *domain.Hosts
	HostMetrics      *domain.Metrics
	HostMetricValues *domain.MetricValues
}

func (repo *HostRepository) List() (*domain.Hosts, error) {
	return repo.Hosts, nil
}

func (repo *HostRepository) Save(host *domain.Host) (*domain.HostID, error) {
	repo.Hosts.Hosts = append(repo.Hosts.Hosts, *host)
	return &domain.HostID{ID: host.ID}, nil
}

func (repo *HostRepository) Host(hostID string) (*domain.Host, error) {
	for i := range repo.Hosts.Hosts {
		if repo.Hosts.Hosts[i].ID == hostID {
			return &repo.Hosts.Hosts[i], nil
		}
	}

	return nil, fmt.Errorf("host not found")
}

func (repo *HostRepository) Exists(hostID string) bool {
	for i := range repo.Hosts.Hosts {
		if repo.Hosts.Hosts[i].ID == hostID {
			return true
		}
	}

	return false
}

func (repo *HostRepository) Status(hostID, status string) (*domain.Success, error) {
	for i := range repo.Hosts.Hosts {
		if repo.Hosts.Hosts[i].ID == hostID {
			repo.Hosts.Hosts[i].Status = status
			return &domain.Success{Success: true}, nil
		}
	}

	return &domain.Success{Success: false}, fmt.Errorf("host not found")
}

func (repo *HostRepository) SaveRoleFullNames(hostID string, names *domain.RoleFullNames) (*domain.Success, error) {
	for i := range repo.Hosts.Hosts {
		if repo.Hosts.Hosts[i].ID == hostID {
			repo.Hosts.Hosts[i].RoleFullNames = names.Names
			return &domain.Success{Success: true}, nil
		}
	}

	return nil, fmt.Errorf("host not found")
}

func (repo *HostRepository) Retire(hostID string, retire *domain.HostRetire) (*domain.Success, error) {
	hosts := []domain.Host{}
	for i := range repo.Hosts.Hosts {
		if repo.Hosts.Hosts[i].ID == hostID {
			continue
		}
		hosts = append(hosts, repo.Hosts.Hosts[i])
	}
	repo.Hosts.Hosts = hosts

	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) ExistsMetric(hostID, name string) bool {
	for i := range repo.HostMetrics.Metrics {
		metric := repo.HostMetrics.Metrics[i]
		if metric.HostID == hostID && metric.Name == name {
			return true
		}
	}

	return false
}

func (repo *HostRepository) MetricNames(hostID string) (*domain.MetricNames, error) {
	names := []string{}
	for i := range repo.HostMetrics.Metrics {
		if repo.HostMetrics.Metrics[i].HostID == hostID {
			names = append(names, repo.HostMetrics.Metrics[i].Name)
		}
	}

	return &domain.MetricNames{Names: names}, nil
}

func (repo *HostRepository) MetricValues(hostID, name string, from, to int) (*domain.MetricValues, error) {
	metrics := []domain.MetricValue{}

	for i := range repo.HostMetricValues.Metrics {
		if repo.HostMetricValues.Metrics[i].HostID != hostID {
			continue
		}
		if repo.HostMetricValues.Metrics[i].Name != name {
			continue
		}
		if from > repo.HostMetricValues.Metrics[i].Time {
			continue
		}
		if repo.HostMetricValues.Metrics[i].Time > to {
			continue
		}

		metrics = append(metrics, repo.HostMetricValues.Metrics[i])
	}

	return &domain.MetricValues{Metrics: metrics}, nil
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

func (repo *HostRepository) MetadataList(hostID string) (*domain.HostMetadata, error) {
	return &domain.HostMetadata{}, nil
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
