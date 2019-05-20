package memory

import (
	"fmt"
	"strings"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	Hosts                  *domain.Hosts
	HostMetadata           []domain.HostMetadata
	HostMetricValues       *domain.MetricValues
	HostMetricValuesLatest map[string]map[string]float64
}

func NewHostRepository() *HostRepository {
	return &HostRepository{
		Hosts:                  &domain.Hosts{Hosts: []domain.Host{}},
		HostMetadata:           []domain.HostMetadata{},
		HostMetricValues:       &domain.MetricValues{},
		HostMetricValuesLatest: make(map[string]map[string]float64),
	}
}

func (repo *HostRepository) List() (*domain.Hosts, error) {
	return repo.Hosts, nil
}

func (repo *HostRepository) Save(host *domain.Host) (*domain.HostID, error) {
	roles := make(map[string][]string)
	for j := range host.RoleFullNames {
		spl := strings.Split(host.RoleFullNames[j], ":")
		roles[spl[0]] = append(roles[spl[0]], spl[1])
	}
	host.Roles = roles

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
	for i := range repo.HostMetricValues.Metrics {
		metric := repo.HostMetricValues.Metrics[i]
		if metric.HostID == hostID && metric.Name == name {
			return true
		}
	}

	return false
}

func (repo *HostRepository) MetricNames(hostID string) (*domain.MetricNames, error) {
	nmap := make(map[string]bool)
	for i := range repo.HostMetricValues.Metrics {
		metric := repo.HostMetricValues.Metrics[i]
		if metric.HostID == hostID {
			nmap[metric.Name] = true
		}
	}

	names := []string{}
	for k := range nmap {
		names = append(names, k)
	}

	return &domain.MetricNames{Names: names}, nil
}

func (repo *HostRepository) MetricValues(hostID, name string, from, to int64) (*domain.MetricValues, error) {
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

func (repo *HostRepository) MetricValuesLatest(hostID, name []string) (*domain.TSDBLatest, error) {
	latest := make(map[string]map[string]float64)
	for i := range hostID {
		if _, ok := repo.HostMetricValuesLatest[hostID[i]]; !ok {
			continue
		}

		if _, ok := latest[hostID[i]]; !ok {
			latest[hostID[i]] = make(map[string]float64)
		}

		for j := range name {
			if v, ok := repo.HostMetricValuesLatest[hostID[i]][name[j]]; ok {
				latest[hostID[i]][name[j]] = v
			}
		}
	}

	return &domain.TSDBLatest{
		TSDBLatest: latest,
	}, nil
}

func (repo *HostRepository) SaveMetricValues(values []domain.MetricValue) (*domain.Success, error) {
	repo.HostMetricValues.Metrics = append(repo.HostMetricValues.Metrics, values...)

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
