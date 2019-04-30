package database

import (
	"fmt"

	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type ServiceMetricRepository struct {
	Internal domain.ServiceMetricValues
}

func NewServiceMetricRepository() *ServiceMetricRepository {
	return &ServiceMetricRepository{
		Internal: domain.ServiceMetricValues{},
	}
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

	return list, fmt.Errorf("service metric not found")
}

func (repo *ServiceMetricRepository) Save(v domain.ServiceMetricValue) error {
	repo.Internal = append(repo.Internal, v)
	return nil
}
