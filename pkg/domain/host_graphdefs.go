package domain

type CustomHostMetricDefs []CustomHostMetricDef

type CustomHostMetricDef struct {
	Name        string      `json:"name"`
	DisplayName string      `json:"displayName,omitempty"`
	Unit        string      `json:"unit,omitempty"`
	Metrics     HostMetrics `json:"metrics"`
}
