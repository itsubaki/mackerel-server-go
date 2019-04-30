package controllers

import (
	"github.com/itsubaki/mackerel-api/pkg/domain"
)

type GetUserInput struct {
}

type GetUserOutput struct {
	Users domain.Users `json:"users"`
}

type DeleteUserInput struct {
	UserID string `json:"-"`
}

type DeleteUserOutput domain.User

func (m *Mackerel) GetUser(in *GetUserInput) (*GetUserOutput, error) {
	return &GetUserOutput{}, nil
}

func (m *Mackerel) DeleteUser(in *DeleteUserInput) (*DeleteUserOutput, error) {
	return &DeleteUserOutput{}, nil
}
