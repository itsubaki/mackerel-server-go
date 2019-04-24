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

type PostServiceInput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type PostServiceOutput struct {
	Status int `json:"-"`
	Service
}

type DeleteServiceInput struct {
	ServiceName string
}

type DeleteServiceOutput struct {
	Status int `json:"-"`
	Service
}
