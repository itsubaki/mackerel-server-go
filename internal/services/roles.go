package services

type Role struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type GetRolesInput struct {
	ServiceName string
}

type GetRolesOutput struct {
	Status int    `json:"-"`
	Roles  []Role `json:"roles"`
}

type PostRolesInput struct {
	ServiceName string `json:"serviceName"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}

type PostRolesOutput struct {
	Status int    `json:"-"`
	Name   string `json:"name"`
	Memo   string `json:"memo"`
}

type DeleteRolesInput struct {
	ServiceName string `json:"serviceName"`
	RoleName    string `json:"name"`
}

type DeleteRolesOutput struct {
	Status int    `json:"-"`
	Name   string `json:"name"`
	Memo   string `json:"memo"`
}
