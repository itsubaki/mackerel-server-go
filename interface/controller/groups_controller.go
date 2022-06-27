package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
)

type NotificationGroupController struct {
	Interactor *usecase.NotificationGroupInteractor
}

func NewNotificationGroupController(handler database.SQLHandler) *NotificationGroupController {
	return &NotificationGroupController{
		Interactor: &usecase.NotificationGroupInteractor{
			NotificationGroupRepository: database.NewNotificationGroupRepository(handler),
		},
	}
}

func (cntr *NotificationGroupController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *NotificationGroupController) Save(c Context) {
	var in domain.NotificationGroup
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

func (cntr *NotificationGroupController) Update(c Context) {
	var in domain.NotificationGroup
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("notificationGroupId")

	out, err := cntr.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *NotificationGroupController) Delete(c Context) {
	out, err := cntr.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("notificationGroupId"),
	)

	doResponse(c, out, err)
}
