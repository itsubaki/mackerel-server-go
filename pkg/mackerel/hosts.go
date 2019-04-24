package mackerel

type GetHostOutput struct {
	Host Host `json:"host"`
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
