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

func (m *Mackerel) PostServices(in *PostServicesInput) *PostServicesOutput {
	return &PostServicesOutput{}
}

func (m *Mackerel) DeleteServices(in *DeleteServicesInput) *DeleteServicesOutput {
	return &DeleteServicesOutput{}
}

func (m *Mackerel) GetRoles(in *GetRolesInput) *GetRolesOutput {
	return &GetRolesOutput{}
}

func (m *Mackerel) PostRoles(in *PostRolesInput) *PostRolesOutput {
	return &PostRolesOutput{}
}

func (m *Mackerel) DeleteRoles(in *DeleteRolesInput) *DeleteRolesOutput {
	return &DeleteRolesOutput{}
}

func (m *Mackerel) GetMetricNames(in *GetMetricNamesInput) *GetMetricNamesOutput {
	return &GetMetricNamesOutput{}
}

func (m *Mackerel) GetHosts(in *GetHostsInput) *GetHostsOutput {
	return &GetHostsOutput{}
}
