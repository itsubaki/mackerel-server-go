package memory

import (
	"fmt"
	"strings"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	Hosts                  *domain.Hosts
	HostMetadata           []domain.HostMetadata
	HostMetrics            *domain.Metrics
	HostMetricValues       *domain.MetricValues
	HostMetricValuesLatest map[string]map[string]float64
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
			roles := make(map[string][]string)
			for j := range names.Names {
				spl := strings.Split(names.Names[j], ":")
				roles[spl[0]] = append(roles[spl[0]], spl[1])
			}

			repo.Hosts.Hosts[i].Roles = roles
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
	return &domain.TSDBLatest{
		TSDBLatest: repo.HostMetricValuesLatest,
	}, nil
}

func (repo *HostRepository) SaveMetricValues(values []domain.MetricValue) (*domain.Success, error) {
	repo.HostMetricValues.Metrics = append(repo.HostMetricValues.Metrics, values...)

	for i := range values {
		repo.HostMetrics.Metrics = append(repo.HostMetrics.Metrics, domain.Metric{
			HostID: values[i].HostID,
			Name:   values[i].Name,
		})
	}

	for i := range values {
		if _, ok := repo.HostMetricValuesLatest[values[i].HostID]; !ok {
			repo.HostMetricValuesLatest[values[i].HostID] = make(map[string]float64)
		}

		repo.HostMetricValuesLatest[values[i].HostID][values[i].Name] = values[i].Value
	}

	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) ExistsMetadata(hostID, namespace string) bool {
	for _, m := range repo.HostMetadata {
		if m.HostID == hostID && m.Namespace == namespace {
			return true
		}
	}

	return false
}

func (repo *HostRepository) MetadataList(hostID string) (*domain.HostMetadataList, error) {
	names := []domain.Namespace{}
	for i := range repo.HostMetadata {
		names = append(names, domain.Namespace{Namespace: repo.HostMetadata[i].Namespace})
	}

	return &domain.HostMetadataList{
		Metadata: names,
	}, nil
}

func (repo *HostRepository) Metadata(hostID, namespace string) (interface{}, error) {
	for i := range repo.HostMetadata {
		if repo.HostMetadata[i].HostID == hostID && repo.HostMetadata[i].Namespace == namespace {
			return repo.HostMetadata[i].Metadata, nil
		}
	}
	return nil, fmt.Errorf("metadata not found")
}

func (repo *HostRepository) SaveMetadata(hostID, namespace string, metadata interface{}) (*domain.Success, error) {
	for i := range repo.HostMetadata {
		if repo.HostMetadata[i].HostID == hostID && repo.HostMetadata[i].Namespace == namespace {
			repo.HostMetadata[i].Metadata = metadata
			return &domain.Success{Success: true}, nil
		}
	}

	repo.HostMetadata = append(repo.HostMetadata, domain.HostMetadata{
		HostID:    hostID,
		Namespace: namespace,
		Metadata:  metadata,
	})

	return &domain.Success{Success: true}, nil
}

func (repo *HostRepository) DeleteMetadata(hostID, namespace string) (*domain.Success, error) {
	list := []domain.HostMetadata{}
	for i := range repo.HostMetadata {
		if repo.HostMetadata[i].HostID == hostID && repo.HostMetadata[i].Namespace == namespace {
			continue
		}
		list = append(list, repo.HostMetadata[i])
	}
	repo.HostMetadata = list

	return &domain.Success{Success: true}, nil
}
