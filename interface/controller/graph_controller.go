package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
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

func (cntr *GraphController) SaveDef(c Context) {
	var in []domain.GraphDef
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveDef(
		c.GetString("org_id"),
		in,
	)

	doResponse(c, out, err)
}

func (cntr *GraphController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *GraphController) Save(c Context) {
	var in domain.GraphAnnotation
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

func (cntr *GraphController) Update(c Context) {
	var in domain.GraphAnnotation
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.Update(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *GraphController) Delete(c Context) {
	out, err := cntr.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("annotationId"),
	)

	doResponse(c, out, err)
}
