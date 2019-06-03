package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type ChannelController struct {
	Interactor *usecase.ChannelInteractor
}

func NewChannelController(handler database.SQLHandler) *ChannelController {
	return &ChannelController{
		Interactor: &usecase.ChannelInteractor{
			ChannelRepository: database.NewChannelRepository(handler),
		},
	}
}

func (s *ChannelController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *ChannelController) Save(c Context) {
	var in domain.Channel
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *ChannelController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("channelId"),
	)

	doResponse(c, out, err)
}
