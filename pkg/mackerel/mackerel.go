package mackerel

func New() (*Mackerel, error) {
	return &Mackerel{}, nil
}

type Mackerel struct{}

func (m *Mackerel) GetServices() *GetServicesOutput {
	return &GetServicesOutput{
		Services: []Service{},
	}
}

func (m *Mackerel) PostService(in *PostServiceInput) *PostServiceOutput {
	return &PostServiceOutput{}
}

func (m *Mackerel) DeleteService(in *DeleteServiceInput) *DeleteServiceOutput {
	return &DeleteServiceOutput{}
}

func (m *Mackerel) GetRoles(in *GetRolesInput) *GetRolesOutput {
	return &GetRolesOutput{}
}

func (m *Mackerel) PostRole(in *PostRoleInput) *PostRoleOutput {
	return &PostRoleOutput{}
}

func (m *Mackerel) DeleteRole(in *DeleteRoleInput) *DeleteRoleOutput {
	return &DeleteRoleOutput{}
}

func (m *Mackerel) GetMetricNames(in *GetMetricNamesInput) *GetMetricNamesOutput {
	return &GetMetricNamesOutput{}
}

func (m *Mackerel) GetHosts() *GetHostsOutput {
	return &GetHostsOutput{}
}

func (m *Mackerel) GetHost(in *GetHostInput) *GetHostOutput {
	return &GetHostOutput{}
}
