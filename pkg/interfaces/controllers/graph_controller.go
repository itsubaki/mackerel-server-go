package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/memory"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type GraphController struct {
	Interactor *usecase.GraphInteractor
}

func NewGraphController(handler database.DB) *GraphController {
	return &GraphController{
		Interactor: &usecase.GraphInteractor{
			GraphRepository: memory.NewGraphRepository(),
		},
	}
}

func (s *GraphController) Save(c Context) {
	var in []domain.GraphDef
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(in)
	doResponse(c, out, err)
}
