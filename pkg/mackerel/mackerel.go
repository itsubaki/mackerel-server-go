package mackerel

import (
	"regexp"
)

func New() (*Mackerel, error) {
	return &Mackerel{
		ServiceNameRule:         regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
		RoleNameRule:            regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
		MetricNameRule:          regexp.MustCompile(`[a-zA-Z0-9._-]+`),
		ServiceRepository:       NewServiceRepository(),
		RoleRepository:          NewRoleRepository(),
		ServiceMetricRepository: NewServiceMetricRepository(),
		HostRepository:          NewHostRepository(),
		HostMetricRepository:    NewHostMetricRepository(),
	}, nil
}

type Mackerel struct {
	ServiceNameRule         *regexp.Regexp
	RoleNameRule            *regexp.Regexp
	MetricNameRule          *regexp.Regexp
	ServiceRepository       *ServiceRepository
	RoleRepository          *RoleRepository
	ServiceMetricRepository *ServiceMetricRepository
	HostRepository          *HostRepository
	HostMetricRepository    *HostMetricRepository
}

func (m *Mackerel) GetServices(in *GetServicesInput) (*GetServicesOutput, error) {
	list, err := m.ServiceRepository.FindAll()
	return &GetServicesOutput{Services: list}, err
}

func (m *Mackerel) PostService(in *PostServiceInput) (*PostServiceOutput, error) {
	if !m.ServiceNameRule.Match([]byte(in.Name)) {
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
		Name:  in.Name,
		Memo:  in.Memo,
		Roles: []string{},
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
	if !m.ServiceNameRule.Match([]byte(in.Name)) {
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

func (m *Mackerel) GetRoleMetadata(in *GetRoleMetadataInput) (GetRoleMetadataOutput, error) {
	return "", nil
}

func (m *Mackerel) PutRoleMetadata(in *PutRoleMetadataInput) (*PutRoleMetadataOutput, error) {
	return &PutRoleMetadataOutput{}, nil
}

func (m *Mackerel) DeleteRoleMetadata(in *DeleteRoleMetadataInput) (*DeleteRoleMetadataOutput, error) {
	return &DeleteRoleMetadataOutput{}, nil
}

func (m *Mackerel) GetRoleMetadataList(in *GetRoleMetadataListInput) (*GetRoleMetadataListOutput, error) {
	return &GetRoleMetadataListOutput{}, nil
}

func (m *Mackerel) GetServiceMetricNames(in *GetServiceMetricNamesInput) (*GetServiceMetricNamesOutput, error) {
	return &GetServiceMetricNamesOutput{}, nil
}

func (m *Mackerel) PostServiceMetric(in *PostServiceMetricInput) (*PostServiceMetricOutput, error) {
	return &PostServiceMetricOutput{}, nil
}

func (m *Mackerel) GetServiceMetric(in *GetServiceMetricInput) (*GetServiceMetricOutput, error) {
	return &GetServiceMetricOutput{}, nil
}

func (m *Mackerel) GetServiceMetadata(in *GetServiceMetadataInput) (GetServiceMetadataOutput, error) {
	return "", nil
}

func (m *Mackerel) PutServiceMetadata(in *PutServiceMetadataInput) (*PutServiceMetadataOutput, error) {
	return &PutServiceMetadataOutput{}, nil
}

func (m *Mackerel) DeleteServiceMetadata(in *DeleteServiceMetadataInput) (*DeleteServiceMetadataOutput, error) {
	return &DeleteServiceMetadataOutput{}, nil
}

func (m *Mackerel) GetServiceMetadataList(in *GetServiceMetadataListInput) (*GetServiceMetadataListOutput, error) {
	return &GetServiceMetadataListOutput{}, nil
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

func (m *Mackerel) PostHostMetric(in *PostHostMetricInput) (*PostHostMetricOutput, error) {
	return &PostHostMetricOutput{}, nil
}

func (m *Mackerel) GetHostMetric(in *GetHostMetricInput) (*GetHostMetricOutput, error) {
	return &GetHostMetricOutput{}, nil
}

func (m *Mackerel) GetHostMetricLatest(in *GetHostMetricLatestInput) (*GetHostMetricLatestOutput, error) {
	return &GetHostMetricLatestOutput{}, nil
}

func (m *Mackerel) PostCustomHostMetricDef(in *PostCustomHostMetricDefInput) (*PostCustomHostMetricDefOutput, error) {
	return &PostCustomHostMetricDefOutput{}, nil
}

func (m *Mackerel) GetHostMetricNames(in *GetHostMetricNamesInput) (*GetHostMetricNamesOutput, error) {
	return &GetHostMetricNamesOutput{}, nil
}

func (m *Mackerel) GetHostMetadata(in *GetHostMetadataInput) (GetHostMetadataOutput, error) {
	return "", nil
}

func (m *Mackerel) PutHostMetadata(in *PutHostMetadataInput) (*PutHostMetadataOutput, error) {
	return &PutHostMetadataOutput{}, nil
}

func (m *Mackerel) DeleteHostMetadata(in *DeleteHostMetadataInput) (*DeleteHostMetadataOutput, error) {
	return &DeleteHostMetadataOutput{}, nil
}

func (m *Mackerel) PostCheckReport(in *PostCheckReportInput) (*PostCheckReportOutput, error) {
	return &PostCheckReportOutput{}, nil
}

func (m *Mackerel) GetAlert(in *GetAlertInput) (*GetAlertOutput, error) {
	return &GetAlertOutput{}, nil
}

func (m *Mackerel) PostAlert(in *PostAlertInput) (*PostAlertOutput, error) {
	return &PostAlertOutput{}, nil
}

func (m *Mackerel) GetUser(in *GetUserInput) (*GetUserOutput, error) {
	return &GetUserOutput{}, nil
}

func (m *Mackerel) DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error) {
	return &DeleteUserOutput{}, nil
}
