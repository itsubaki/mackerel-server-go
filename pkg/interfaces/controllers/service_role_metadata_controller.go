package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GetRoleMetadataInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
	Namespace   string `json:"-"`
}

type GetRoleMetadataOutput interface{}

type PutRoleMetadataInput struct {
	ServiceName string      `json:"-"`
	RoleName    string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type PutRoleMetadataOutput struct {
	Success bool `json:"success"`
}

type DeleteRoleMetadataInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
	Namespace   string `json:"-"`
}

type DeleteRoleMetadataOutput struct {
	Success bool `json:"success"`
}

type GetRoleMetadataListInput struct {
	ServiceName string `json:"-"`
	RoleName    string `json:"-"`
}

type GetRoleMetadataListOutput struct {
	Metadata domain.ServiceRoleMetadataList `json:"metadata"`
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
