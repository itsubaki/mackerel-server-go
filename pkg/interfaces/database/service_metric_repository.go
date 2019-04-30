package database

import "github.com/itsubaki/mackerel-api/pkg/domain"

type ServiceMetricRepository struct {
	Internal domain.ServiceMetricValues
}

func (repo *ServiceMetricRepository) FindAll() (domain.ServiceMetricValues, error) {
	return repo.Internal, nil
}

func (repo *ServiceMetricRepository) FindBy(serviceName, metricName string, from, to int64) (domain.ServiceMetricValues, error) {
	list := domain.ServiceMetricValues{}
	for i := range repo.Internal {
		if repo.Internal[i].ServiceName != serviceName {
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

func (repo *ServiceMetricRepository) Save(v domain.ServiceMetricValues) error {
	repo.Internal = append(repo.Internal, v...)
	return nil
}
