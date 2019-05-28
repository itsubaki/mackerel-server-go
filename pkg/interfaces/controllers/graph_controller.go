package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type GraphController struct {
	Interactor *usecase.GraphInteractor
}

func NewGraphController(handler database.SQLHandler) *GraphController {
	return &GraphController{
		Interactor: &usecase.GraphInteractor{
			GraphRepository: database.NewGraphRepository(handler),
		},
	}
}

func (s *GraphController) SaveDef(c Context) {
	var in []domain.GraphDef
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveDef(
		c.GetString("org_id"),
		in,
	)

	doResponse(c, out, err)
}
