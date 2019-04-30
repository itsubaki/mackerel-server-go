package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type HostMetricRepository struct {
	Internal domain.HostMetricValues
}

func NewHostMetricRepository() *HostMetricRepository {
	return &HostMetricRepository{
		Internal: domain.HostMetricValues{},
	}
}

func (repo *HostMetricRepository) FindAll() (domain.HostMetricValues, error) {
	return repo.Internal, nil
}

func (repo *HostMetricRepository) Latest(hostID, metricName []string) (domain.HostMetricValues, error) {
	return domain.HostMetricValues{}, nil
}

func (repo *HostMetricRepository) ExistsByName(hostID, metricName string) bool {
	for i := range repo.Internal {
		if repo.Internal[i].HostID == hostID && repo.Internal[i].Name == metricName {
			return true
		}
	}

	return false
}

func (repo *HostMetricRepository) FindByID(hostID string) (domain.HostMetricValues, error) {
	list := domain.HostMetricValues{}
	for i := range repo.Internal {
		if repo.Internal[i].HostID != hostID {
			continue
		}
		list = append(list, repo.Internal[i])
	}

	return list, nil
}

func (repo *HostMetricRepository) FindBy(hostID, metricName string, from, to int64) (domain.HostMetricValues, error) {
	list := domain.HostMetricValues{}
	for i := range repo.Internal {
		if repo.Internal[i].HostID != hostID {
			continue
		}
		if repo.Internal[i].Name != metricName {
			continue
		}
		if from > repo.Internal[i].Time {
			continue
		}
		if repo.Internal[i].Time > to {
			continue
		}

		list = append(list, repo.Internal[i])
	}

	return list, nil
}

func (repo *HostMetricRepository) Save(v domain.HostMetricValues) error {
	repo.Internal = append(repo.Internal, v...)
	return nil
}
