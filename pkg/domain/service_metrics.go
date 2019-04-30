package domain

type ServiceMetricValues []ServiceMetricValue

type ServiceMetricValue struct {
	ServiceName string  `json:"-"`
	Name        string  `json:"name"`
	Time        int64   `json:"time"`
	Value       float64 `json:"value"`
}

func (v ServiceMetricValues) MetricNames() []string {
	nmap := make(map[string]bool)
	for i := range v {
		nmap[v[i].Name] = true
	}

	names := []string{}
	for k := range nmap {
		names = append(names, k)
	}

	return names
}
