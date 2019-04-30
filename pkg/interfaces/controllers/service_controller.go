package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GetServicesInput struct {
}

type GetServicesOutput struct {
	Services domain.Services `json:"services"`
}

type PostServiceInput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type PostServiceOutput domain.Service

type DeleteServiceInput struct {
	ServiceName string `json:"-"`
}

type DeleteServiceOutput domain.Service

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

	if err := m.ServiceRepository.Save(domain.Service{
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

	r, err := m.ServiceRoleRepository.FindAll(in.ServiceName)
	if err != nil {
		return nil, err
	}

	for i := range r {
		if err := m.ServiceRoleRepository.Delete(r[i].ServiceName, r[i].Name); err != nil {
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
