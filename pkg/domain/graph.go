package domain

type GraphDef struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Metrics     []Metric `json:"metrics"`
}
