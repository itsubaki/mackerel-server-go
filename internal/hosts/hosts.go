package hosts

import "encoding/json"

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

func (o *GetHostOutput) String() string {
	bytea, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}
