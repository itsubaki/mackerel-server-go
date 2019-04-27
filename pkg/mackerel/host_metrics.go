package mackerel

import "fmt"

type PostHostMetricInput struct {
	MetricValue []HostMetricValue `json:"-"`
}

type PostHostMetricOutput struct {
	Success bool `json:"success"`
}

type GetHostMetricInput struct {
	HostID string `json:"-"`
	Name   string `json:"-"`
	From   string `json:"-"`
	To     string `json:"-"`
}

type GetHostMetricOutput struct {
	Metrics []HostMetricValue `json:"metrics"`
}

type GetHostMetricNamesInput struct {
	HostID string `json:"-"`
}

type GetHostMetricNamesOutput struct {
	Name []string `json:"names"`
}

type GetHostMetricLatestInput struct {
	HostID string `json:"-"`
	Name   string `json:"-"`
}

type GetHostMetricLatestOutput struct {
	TSDBLatest map[string]map[string]float64 `json:"tsdbLatest"`
}

type HostMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	IsStacked   bool   `json:"isStacked"`
}

type HostMetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}

type HostMetricRepository struct {
	Internal []HostMetricValue
}

func NewHostMetricRepository() *HostMetricRepository {
	return &HostMetricRepository{
		Internal: []HostMetricValue{},
	}
}

func (repo *HostMetricRepository) Latest(hostID, metricName []string) ([]HostMetricValue, error) {
	return []HostMetricValue{}, nil
}

func (repo *HostMetricRepository) ExistsByName(hostID, metricName string) bool {
	for i := range repo.Internal {
		if repo.Internal[i].HostID == hostID && repo.Internal[i].Name == metricName {
			return true
		}
	}

	return false
}

func (repo *HostMetricRepository) FindBy(hostID, metricName string, from, to int64) ([]HostMetricValue, error) {
	list := []HostMetricValue{}
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

	return list, fmt.Errorf("host metric not found")
}

func (repo *HostMetricRepository) Save(v HostMetricValue) error {
	repo.Internal = append(repo.Internal, v)
	return nil
}
