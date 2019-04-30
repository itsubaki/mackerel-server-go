package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type PostHostInput struct {
	domain.Host
}

type PostHostOutput struct {
	ID string `json:"id"`
}

type GetHostInput struct {
	HostID string `json:"-"`
}

type GetHostOutput struct {
	Host domain.Host `json:"host"`
}

type PutHostInput struct {
	HostID string `json:"-"`
	domain.Host
}

type PutHostOutput struct {
	ID string `json:"id"`
}

type PostHostStatusInput struct {
	HostID string `json:"-"`
	Status string `json:"status"` // standby, working, maintenance, poweroff
}

type PostHostStatusOutput struct {
	Success bool `json:"success"`
}

type PutHostRoleFullNamesInput struct {
	HostID        string   `json:"-"`
	RollFullNames []string `json:"roleFullnames"`
}

type PutHostRoleFullNamesOutput struct {
	Success bool `json:"success"`
}

type PostHostRetiredInput struct {
	HostID string `json:"-"`
}

type PostHostRetiredOutput struct {
	Success bool `json:"success"`
}

type GetHostsInput struct {
	ServiceName      string   `json:"-"`
	RoleName         []string `json:"-"`
	Name             string   `json:"-"`
	Status           string   `json:"-"`
	CustomIdentifier string   `json:"-"`
}

type GetHostsOutput struct {
	Host domain.Hosts `json:"hosts"`
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
