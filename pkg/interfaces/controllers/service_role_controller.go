package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GetServiceRolesInput struct {
	ServiceName string `json:"-"`
}

type GetServiceRolesOutput struct {
	Roles domain.ServiceRoles `json:"roles"`
}

type PostServiceRoleInput struct {
	ServiceName string `json:"serviceName"`
	Name        string `json:"name"`
	Memo        string `json:"memo"`
}

type PostServiceRoleOutput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

type DeleteServiceRoleInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"name"`
}

type DeleteServiceRoleOutput struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

func (m *Mackerel) GetServiceRoles(in *GetServiceRolesInput) (*GetServiceRolesOutput, error) {
	list, err := m.ServiceRoleRepository.FindAll(in.ServiceName)
	if err != nil {
		return nil, &ServiceNotFound{}
	}

	return &GetServiceRolesOutput{Roles: list}, nil
}

func (m *Mackerel) PostServiceRole(in *PostServiceRoleInput) (*PostServiceRoleOutput, error) {
	if !m.ServiceNameRule.Match([]byte(in.Name)) {
		return nil, &InvalidRoleName{}
	}

	if m.ServiceRoleRepository.ExistsByName(in.ServiceName, in.Name) {
		return nil, &InvalidRoleName{}
	}

	m.ServiceRoleRepository.Save(domain.ServiceRole{
		ServiceName: in.ServiceName,
		Name:        in.Name,
		Memo:        in.Memo,
	})

	return &PostServiceRoleOutput{
		Name: in.Name,
		Memo: in.Memo,
	}, nil
}

func (m *Mackerel) DeleteServiceRole(in *DeleteServiceRoleInput) (*DeleteServiceRoleOutput, error) {
	r, err := m.ServiceRoleRepository.FindByName(in.ServiceName, in.RoleName)
	if err != nil {
		return nil, &RoleNotFound{}
	}

	if err := m.ServiceRoleRepository.Delete(in.ServiceName, in.RoleName); err != nil {
		return nil, err
	}

	return &DeleteServiceRoleOutput{
		Name: r.Name,
		Memo: r.Memo,
	}, nil
}
