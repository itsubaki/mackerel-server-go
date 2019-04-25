package mackerel

type GetRolesInput struct {
	ServiceName string
}

type GetRolesOutput struct {
	Roles []Role `json:"roles"`
}

type PostRoleInput struct {
	ServiceName string `json:"serviceName"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}

type PostRoleOutput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type DeleteRoleInput struct {
	ServiceName string `json:"serviceName"`
	RoleName    string `json:"name"`
}

type DeleteRoleOutput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type Role struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}
