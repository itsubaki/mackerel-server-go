package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
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

func (s *GraphController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *GraphController) Save(c Context) {
	var in domain.GraphAnnotation
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

func (s *GraphController) Update(c Context) {
	var in domain.GraphAnnotation
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *GraphController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("annotationId"),
	)

	doResponse(c, out, err)
}
