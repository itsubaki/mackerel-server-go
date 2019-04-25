package mackerel

import "regexp"

func New() (*Mackerel, error) {
	return &Mackerel{
		ServiceRepository:     NewServiceRepository(),
		RoleRepository:        NewRoleRepository(),
		MetricValueRepository: NewMetricValueRepository(),
		HostRepository:        NewHostRepository(),
	}, nil
}

type Mackerel struct {
	ServiceRepository     *ServiceRepository
	RoleRepository        *RoleRepository
	MetricValueRepository *MetricValueRepository
	HostRepository        *HostRepository
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

	r, err := m.RoleRepository.FindAll(in.ServiceName)
	if err != nil {
		return nil, err
	}

	for i := range r {
		if err := m.RoleRepository.Delete(r[i].ServiceName, r[i].Name); err != nil {
			return nil, err
		}
	}

	if err := m.ServiceRepository.Delete(in.ServiceName); err != nil {
		return nil, err
	}

	roles := []string{}
	for i := range r {
		roles = append(roles, r[i].Name)
	}

	return &DeleteServiceOutput{
		Name:  s.Name,
		Memo:  s.Memo,
		Roles: roles,
	}, nil
}

func (m *Mackerel) GetRoles(in *GetRolesInput) (*GetRolesOutput, error) {
	list, err := m.RoleRepository.FindAll(in.ServiceName)
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	return &GetRolesOutput{Roles: list}, nil
}

func (m *Mackerel) PostRole(in *PostRoleInput) (*PostRoleOutput, error) {
	if !regexp.MustCompile(`^[a-zA-Z]{1,1}[a-zA-Z0-9_-]{1,62}`).Match([]byte(in.Name)) {
		return nil, &InvalidRoleName{}
	}

	if m.RoleRepository.Exist(in.ServiceName, in.Name) {
		return nil, &InvalidRoleName{}
	}

	m.RoleRepository.Insert(Role{
		ServiceName: in.ServiceName,
		Name:        in.Name,
		Memo:        in.Memo,
	})

	return &PostRoleOutput{
		Name: in.Name,
		Memo: in.Memo,
	}, nil
}

func (m *Mackerel) DeleteRole(in *DeleteRoleInput) (*DeleteRoleOutput, error) {
	r, err := m.RoleRepository.Find(in.ServiceName, in.RoleName)
	if err != nil {
		return nil, &RoleNotFound{}
	}

	if err := m.RoleRepository.Delete(in.ServiceName, in.RoleName); err != nil {
		return nil, err
	}

	return &DeleteRoleOutput{
		Name: r.Name,
		Memo: r.Memo,
	}, nil
}

func (m *Mackerel) GetMetricNames(in *GetMetricNamesInput) (*GetMetricNamesOutput, error) {
	return &GetMetricNamesOutput{}, nil
}

func (m *Mackerel) PostHost(in *PostHostInput) (*PostHostOutput, error) {
	return &PostHostOutput{}, nil
}

func (m *Mackerel) GetHosts() (*GetHostsOutput, error) {
	return &GetHostsOutput{}, nil
}

func (m *Mackerel) GetHost(in *GetHostInput) (*GetHostOutput, error) {
	return &GetHostOutput{}, nil
}
