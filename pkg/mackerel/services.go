package mackerel

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
