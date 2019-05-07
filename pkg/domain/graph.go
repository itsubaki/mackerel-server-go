package domain

type GraphDef struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Metrics     []Metric `json:"metrics"`
}

type GraphAnnotations struct {
	GraphAnnotations []GraphAnnotation `json:"graphAnnotations"`
}

type GraphAnnotation struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	From        int64    `json:"from"`
	To          int64    `json:"to"`
	Service     string   `json:"service"`
	Roles       []string `json:"roles"`
}
