package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
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

func (s *NotificationGroupController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *NotificationGroupController) Save(c Context) {
	var in domain.NotificationGroup
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

func (s *NotificationGroupController) Update(c Context) {
	var in domain.NotificationGroup
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("notificationGroupId")

	out, err := s.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *NotificationGroupController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("notificationGroupId"),
	)

	doResponse(c, out, err)
}
