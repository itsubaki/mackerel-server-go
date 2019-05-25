package domain

type GraphDef struct {
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

type GraphAnnotations struct {
	GraphAnnotations []GraphAnnotation `json:"graphAnnotations"`
}

type GraphAnnotation struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	From        int64    `json:"from"`
	To          int64    `json:"to"`
	Service     string   `json:"service"`
	Roles       []string `json:"roles,omitempty"`
}
