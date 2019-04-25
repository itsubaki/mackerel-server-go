package mackerel

import "regexp"

func New() (*Mackerel, error) {
	return &Mackerel{
		ServiceRepository: NewServiceRepository(),
	}, nil
}

type Mackerel struct {
	ServiceRepository *ServiceRepository
}

func (m *Mackerel) GetServices() (*GetServicesOutput, error) {
	list, err := m.ServiceRepository.FindAll()
	return &GetServicesOutput{Services: list}, err
}

func (m *Mackerel) PostService(in *PostServiceInput) (*PostServiceOutput, error) {
	if !regexp.MustCompile(`^[a-zA-Z]{1,1}[a-zA-Z0-9_-]{1,62}`).Match([]byte(in.Name)) {
		return nil, &InvalidServiceName{}
	}

	if m.ServiceRepository.Exist(in.Name) {
		return nil, &InvalidServiceName{}
	}

	if err := m.ServiceRepository.Insert(Service{
		Name:  in.Name,
		Memo:  in.Memo,
		Roles: []string{},
	}); err != nil {
		return nil, err
	}

	return &PostServiceOutput{
		Name:  in.Name,
		Memo:  in.Memo,
		Roles: []string{},
	}, nil
}

func (m *Mackerel) DeleteService(in *DeleteServiceInput) (*DeleteServiceOutput, error) {
	s, err := m.ServiceRepository.Find(Service{Name: in.ServiceName})
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	if err := m.ServiceRepository.Delete(in.ServiceName); err != nil {
		return nil, err
	}

	return &DeleteServiceOutput{
		Name:  s.Name,
		Memo:  s.Memo,
		Roles: s.Roles,
	}, nil
}

func (m *Mackerel) GetRoles(in *GetRolesInput) (*GetRolesOutput, error) {
	return &GetRolesOutput{}, nil
}

func (m *Mackerel) PostRole(in *PostRoleInput) (*PostRoleOutput, error) {
	return &PostRoleOutput{}, nil
}

func (m *Mackerel) DeleteRole(in *DeleteRoleInput) (*DeleteRoleOutput, error) {
	return &DeleteRoleOutput{}, nil
}

func (m *Mackerel) GetMetricNames(in *GetMetricNamesInput) (*GetMetricNamesOutput, error) {
	return &GetMetricNamesOutput{}, nil
}

func (m *Mackerel) GetHosts() (*GetHostsOutput, error) {
	return &GetHostsOutput{}, nil
}

func (m *Mackerel) GetHost(in *GetHostInput) (*GetHostOutput, error) {
	return &GetHostOutput{}, nil
}
