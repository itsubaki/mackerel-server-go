package services

import "encoding/json"

type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

type GetServicesOutput struct {
	Services []Service `json:"services"`
}

func (o *GetServicesOutput) String() string {
	bytea, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}
