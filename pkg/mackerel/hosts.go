package mackerel

type GetHostInput struct {
	HostID string
}

type GetHostOutput struct {
	Status int  `json:"-"`
	Host   Host `json:"host"`
}

type GetHostsOutput struct {
	Status int    `json:"-"`
	Host   []Host `json:"hosts"`
}

type Host struct {
	CreatedAt   int64    `json:"createdAt"`
	ID          string   `json:"id"`
	Status      string   `json:"status"`
	Memo        string   `json:"memo"`
	Roles       []string `json:"roles"`
	Interfaces  []string `json:"interfaces"`
	IsRetired   bool     `json:"isRetired"`
	DisplayName string   `json:"displayName"`
	Meta        string   `json:"meta"`
}
