package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type HostRepository struct {
	SQLHandler       SQLHandler
	Hosts            domain.Hosts
	HostMetadata     domain.HostMetadataList
	HostMetricValues domain.HostMetricValues
	CustomGraphDefs  domain.CustomGraphDefs
}

func (repo *HostRepository) SaveCustomGraphDefs(v domain.CustomGraphDefs) error {
	repo.CustomGraphDefs = append(repo.CustomGraphDefs, v...)
	return nil
}

func (repo *HostRepository) MetricNames() ([]string, error) {
	return []string{}, nil
}

func (repo *HostRepository) MetricValuesLatest(hostID, metricName []string) (domain.HostMetricValues, error) {
	return domain.HostMetricValues{}, nil
}

func (repo *HostRepository) MetricValues(hostID, metricName string, from, to int64) (domain.HostMetricValues, error) {
	list := domain.HostMetricValues{}
	for i := range repo.HostMetricValues {
		if repo.HostMetricValues[i].HostID != hostID {
			continue
		}
		if repo.HostMetricValues[i].Name != metricName {
			continue
		}
		if from > repo.HostMetricValues[i].Time {
			continue
		}
		if repo.HostMetricValues[i].Time > to {
			continue
		}

		list = append(list, repo.HostMetricValues[i])
	}

	return list, nil
}

func (repo *HostRepository) SaveMetricValues(v domain.HostMetricValues) error {
	repo.HostMetricValues = append(repo.HostMetricValues, v...)
	return nil
}

func (repo *HostRepository) Exists(hostName string) bool {
	for i := range repo.Hosts {
		if repo.Hosts[i].Name == hostName {
			return true
		}
	}

	return false
}

func (repo *HostRepository) FindByID(hostID string) (domain.Host, error) {
	for i := range repo.Hosts {
		if repo.Hosts[i].ID == hostID {
			return repo.Hosts[i], nil
		}
	}

	return domain.Host{}, fmt.Errorf("host not found")
}

func (repo *HostRepository) FindByName(hostName string) (domain.Host, error) {
	for i := range repo.Hosts {
		if repo.Hosts[i].Name == hostName {
			return repo.Hosts[i], nil
		}
	}

	return domain.Host{}, fmt.Errorf("host not found")
}

func (repo *HostRepository) FindAll() (domain.Hosts, error) {
	return repo.Hosts, nil
}

func (repo *HostRepository) Save(host domain.Host) error {
	repo.Hosts = append(repo.Hosts, host)
	return nil
}

func (repo *HostRepository) DeleteByID(hostID string) error {
	list := domain.Hosts{}
	for i := range repo.Hosts {
		if repo.Hosts[i].ID != hostID {
			list = append(list, repo.Hosts[i])
		}
	}
	repo.Hosts = list
	return nil
}

func (repo *HostRepository) DeleteByName(hostName string) error {
	list := domain.Hosts{}
	for i := range repo.Hosts {
		if repo.Hosts[i].Name != hostName {
			list = append(list, repo.Hosts[i])
		}
	}
	repo.Hosts = list
	return nil
}
