package mackerel

type PostCustomHostMetricDefInput struct {
	CustomHostMetricDef []CustomHostMetricDef `json:"-"`
}

type PostCustomHostMetricDefOutput struct {
	Success bool `json:"success"`
}

type CustomHostMetricDef struct {
	Name        string       `json:"name"`
	DisplayName string       `json:"displayName,omitempty"`
	Unit        string       `json:"unit,omitempty"`
	Metrics     []HostMetric `json:"metrics"`
}

type CustomGraphDefRepository struct {
	Internal []CustomHostMetricDef
}

func NewCustomHostMetricDefRepository() *CustomGraphDefRepository {
	return &CustomGraphDefRepository{
		Internal: []CustomHostMetricDef{},
	}
}

func (repo *CustomGraphDefRepository) Save(v CustomHostMetricDef) error {
	repo.Internal = append(repo.Internal, v)
	return nil
}
