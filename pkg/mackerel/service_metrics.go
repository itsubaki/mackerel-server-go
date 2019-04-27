package mackerel

import "fmt"

type PostServiceMetricInput struct {
	ServiceName        string               `json:"-"`
	ServiceMetricValue []ServiceMetricValue `json:"-"`
}

type PostServiceMetricOutput struct {
	Success bool `json:"success"`
}

type GetServiceMetricInput struct {
	ServiceName string `json:"-"`
	Name        string `json:"-"`
	From        string `json:"-"`
	To          string `json:"-"`
}

type GetServiceMetricOutput struct {
	Metrics []ServiceMetricValue `json:"metrics"`
}

type GetServiceMetricNamesInput struct {
	ServiceName string `json:"-"`
}

type GetServiceMetricNamesOutput struct {
	Name []string `json:"names"`
}

type ServiceMetricValue struct {
	ServiceName string  `json:"-"`
	Name        string  `json:"name"`
	Time        int64   `json:"time"`
	Value       float64 `json:"value"`
}

type ServiceMetricRepository struct {
	Internal []ServiceMetricValue
}

func NewServiceMetricRepository() *ServiceMetricRepository {
	return &ServiceMetricRepository{
		Internal: []ServiceMetricValue{},
	}
}

func (repo *ServiceMetricRepository) FindBy(serviceName, metricName string, from, to int64) ([]ServiceMetricValue, error) {
	list := []ServiceMetricValue{}
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

func (repo *ServiceMetricRepository) Save(v ServiceMetricValue) error {
	repo.Internal = append(repo.Internal, v)
	return nil
}
