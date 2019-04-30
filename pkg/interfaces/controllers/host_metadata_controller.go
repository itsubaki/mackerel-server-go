package controllers

import "github.com/itsubaki/mackerel-api/pkg/domain"

type GetHostMetadataInput struct {
	HostID    string `json:"-"`
	Namespace string `json:"-"`
}

type GetHostMetadataOutput interface{}

type PutHostMetadataInput struct {
	HostID    string      `json:"-"`
	Namespace string      `json:"-"`
	Metadata  interface{} `json:"-"`
}

type PutHostMetadataOutput struct {
	Success bool `json:"success"`
}

type DeleteHostMetadataInput struct {
	HostID    string `json:"-"`
	Namespace string `json:"-"`
}

type DeleteHostMetadataOutput struct {
	Success bool `json:"success"`
}

type GetHostMetadataListInput struct {
	HostID string `json:"-"`
}

type GetHostMetadataListOutput struct {
	Metadata domain.HostMetadataList `json:"metadata"`
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
