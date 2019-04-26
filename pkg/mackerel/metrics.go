package mackerel

type PostHostMetricInput struct {
	MetricValue []MetricValue `json:"-"`
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
	Metrics []MetricValue `json:"metrics"`
}

type GetHostMetricLatestInput struct {
	HostID string `json:"-"`
	Name   string `json:"-"`
}

type GetHostMetricLatestOutput struct {
	TSDBLatest map[string]map[string]float64 `json:"tsdbLatest"`
}

type PostCustomMetricDefInput struct {
	CustomMetricDef []CustomMetricDef `json:"-"`
}

type PostCustomMetricDefOutput struct {
	Success bool `json:"success"`
}

type CustomMetricDef struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Metrics     []Metric `json:"metrics"`
}

type Metric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	IsStacked   bool   `json:"isStacked"`
}

type MetricNames struct {
	Name []string `json:"names"`
}

type MetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}

type MetricValueRepository struct {
	Internal []MetricValue
}

func NewMetricValueRepository() *MetricValueRepository {
	return &MetricValueRepository{
		Internal: []MetricValue{},
	}
}

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
	Metrics []MetricValue `json:"metrics"`
}

type ServiceMetricValue struct {
	Name  string  `json:"name"`
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}
