package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GetServiceMetadataInput struct {
	ServiceName string `json:"-"`
	Namespace   string `json:"-"`
}

type GetServiceMetadataOutput interface{}

type PutServiceMetadataInput struct {
	ServiceName string      `json:"-"`
	Namespace   string      `json:"-"`
	Metadata    interface{} `json:"-"`
}

type PutServiceMetadataOutput struct {
	Success bool `json:"success"`
}

type DeleteServiceMetadataInput struct {
	ServiceName string `json:"-"`
	Namespace   string `json:"-"`
}

type DeleteServiceMetadataOutput struct {
	Success bool `json:"success"`
}

type GetServiceMetadataListInput struct {
	ServiceName string `json:"-"`
}

type GetServiceMetadataListOutput struct {
	Metadata domain.ServiceMetadataList `json:"metadata"`
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
