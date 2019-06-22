package domain

type Dashboards struct {
	Dashboards []Dashboard `json:"dashboards,omitempty"`
}

type Dashboard struct {
	OrgID     string   `json:"-"`
	ID        string   `json:"id,omitempty"`
	Title     string   `json:"title"`
	Memo      string   `json:"memo"`
	URLPath   string   `json:"urlPath"`
	Widgets   []Widget `json:"widgets"`
	CreatedAt int64    `json:"createdAt,omitempty"`
	UpdatedAt int64    `json:"updatedAt,omitempty"`
}

type Widget struct {
	Type   string       `json:"type"` // graph, value, markdown
	Title  string       `json:"title"`
	Layout WidgetLayout `json:"layout"`
}

type WidgetGraph struct {
	Type string `json:"type"` // host, service, role, expression, unknown
}

type WidgetRange struct {
	Type string `json:"type"` // relative, absolute
}

type WidgetMetric struct {
	Type string `json:"type"` // host, service, expression
}

type WidgetLayout struct {
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}
