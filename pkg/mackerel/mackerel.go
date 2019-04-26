package mackerel

import (
	"regexp"
)

func New() (*Mackerel, error) {
	return &Mackerel{
		NameRule:              regexp.MustCompile(`^[a-zA-Z]{1,1}[a-zA-Z0-9_-]{1,62}`),
		ServiceRepository:     NewServiceRepository(),
		RoleRepository:        NewRoleRepository(),
		MetricValueRepository: NewMetricValueRepository(),
		HostRepository:        NewHostRepository(),
	}, nil
}

type Mackerel struct {
	NameRule              *regexp.Regexp
	ServiceRepository     *ServiceRepository
	RoleRepository        *RoleRepository
	MetricValueRepository *MetricValueRepository
	HostRepository        *HostRepository
}

func (m *Mackerel) GetServices(in *GetServicesInput) (*GetServicesOutput, error) {
	list, err := m.ServiceRepository.FindAll()
	return &GetServicesOutput{Services: list}, err
}

func (m *Mackerel) PostService(in *PostServiceInput) (*PostServiceOutput, error) {
	if !m.NameRule.Match([]byte(in.Name)) {
		return nil, &InvalidServiceName{}
	}

	if m.ServiceRepository.ExistsByName(in.Name) {
		return nil, &InvalidServiceName{}
	}

	if err := m.ServiceRepository.Save(Service{
		Name:  in.Name,
		Memo:  in.Memo,
		Roles: []string{},
	}); err != nil {
		return nil, err
	}

	return &PostServiceOutput{
		Service: Service{
			Name:  in.Name,
			Memo:  in.Memo,
			Roles: []string{},
		},
	}, nil
}

func (m *Mackerel) DeleteService(in *DeleteServiceInput) (*DeleteServiceOutput, error) {
	s, err := m.ServiceRepository.FindByName(in.ServiceName)
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
		Service: Service{
			Name:  s.Name,
			Memo:  s.Memo,
			Roles: roles,
		},
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
	if !m.NameRule.Match([]byte(in.Name)) {
		return nil, &InvalidRoleName{}
	}

	if m.RoleRepository.ExistsByName(in.ServiceName, in.Name) {
		return nil, &InvalidRoleName{}
	}

	m.RoleRepository.Save(Role{
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
	r, err := m.RoleRepository.FindByName(in.ServiceName, in.RoleName)
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

func (m *Mackerel) GetHost(in *GetHostInput) (*GetHostOutput, error) {
	return &GetHostOutput{}, nil
}

func (m *Mackerel) PutHost(in *PutHostInput) (*PutHostOutput, error) {
	return &PutHostOutput{}, nil
}

func (m *Mackerel) PostHostStatus(in *PostHostStatusInput) (*PostHostStatusOutput, error) {
	return &PostHostStatusOutput{}, nil
}

func (m *Mackerel) PutHostRoleFullNames(in *PutHostRoleFullNamesInput) (*PutHostRoleFullNamesOutput, error) {
	return &PutHostRoleFullNamesOutput{}, nil
}

func (m *Mackerel) PostHostRetired(in *PostHostRetiredInput) (*PostHostRetiredOutput, error) {
	return &PostHostRetiredOutput{}, nil
}

func (m *Mackerel) GetHosts(in *GetHostsInput) (*GetHostsOutput, error) {
	return &GetHostsOutput{}, nil
}
