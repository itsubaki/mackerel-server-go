package services

import "encoding/json"

type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}

type GetServicesOutput struct {
	Status   int       `json:"-"`
	Services []Service `json:"services"`
}

type PostServicesInput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type PostServicesOutput struct {
	Status int `json:"-"`
	Service
}

type DeleteServicesInput struct {
	ServiceName string
}

type DeleteServicesOutput struct {
	Status int `json:"-"`
	Service
}

func (o *GetServicesOutput) String() string {
	bytea, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (o *PostServicesOutput) String() string {
	bytea, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (o *DeleteServicesOutput) String() string {
	bytea, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}
