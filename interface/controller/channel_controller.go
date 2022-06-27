package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
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

func (cntr *ChannelController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *ChannelController) Save(c Context) {
	var in domain.Channel
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.Save(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *ChannelController) Delete(c Context) {
	out, err := cntr.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("channelId"),
	)

	doResponse(c, out, err)
}
