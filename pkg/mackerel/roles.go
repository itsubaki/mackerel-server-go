package mackerel

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

type PostRoleInput struct {
	ServiceName string `json:"serviceName"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}

type PostRoleOutput struct {
	Status int    `json:"-"`
	Name   string `json:"name"`
	Memo   string `json:"memo"`
}

type DeleteRoleInput struct {
	ServiceName string `json:"serviceName"`
	RoleName    string `json:"name"`
}

type DeleteRoleOutput struct {
	Status int    `json:"-"`
	Name   string `json:"name"`
	Memo   string `json:"memo"`
}
