package controller

import (
	"net/http"
	"strconv"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
)

type HostController struct {
	Interactor *usecase.HostInteractor
}

func NewHostController(handler database.SQLHandler) *HostController {
	return &HostController{
		Interactor: &usecase.HostInteractor{
			HostRepository:       database.NewHostRepository(handler),
			HostMetaRepository:   database.NewHostMetaRepository(handler),
			HostMetricRepository: database.NewHostMetricRepository(handler),
			ServiceRepository:    database.NewServiceRepository(handler),
			RoleRepository:       database.NewRoleRepository(handler),
		},
	}
}

func (cntr *HostController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *HostController) Save(c Context) {
	var in domain.Host
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

func (cntr *HostController) Update(c Context) {
	var in domain.Host
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("hostId")

	out, err := cntr.Interactor.Save(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *HostController) Host(c Context) {
	out, err := cntr.Interactor.Host(
		c.GetString("org_id"),
		c.Param("hostId"),
	)

	doResponse(c, out, err)
}

func (cntr *HostController) Status(c Context) {
	var in domain.HostStatus
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.Status(
		c.GetString("org_id"),
		c.Param("hostId"),
		in.Status,
	)

	doResponse(c, out, err)
}

func (cntr *HostController) RoleFullNames(c Context) {
	var in domain.RoleFullNames
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveRoleFullNames(
		c.GetString("org_id"),
		c.Param("hostId"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *HostController) Retire(c Context) {
	var in domain.HostRetire

	// mkr request dont have empty body.
	//if err := c.BindJSON(&in); err != nil {
	//	c.Status(http.StatusBadRequest)
	//	return
	//}

	out, err := cntr.Interactor.Retire(
		c.GetString("org_id"),
		c.Param("hostId"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *HostController) MetricNames(c Context) {
	out, err := cntr.Interactor.MetricNames(
		c.GetString("org_id"),
		c.Param("hostId"),
	)

	doResponse(c, out, err)
}

func (cntr *HostController) MetricValues(c Context) {
	from, err := strconv.ParseInt(c.Query("from"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	to, err := strconv.ParseInt(c.Query("to"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.MetricValues(
		c.GetString("org_id"),
		c.Param("hostId"),
		c.Query("name"),
		from,
		to,
	)

	doResponse(c, out, err)
}

func (cntr *HostController) MetricValuesLatest(c Context) {
	out, err := cntr.Interactor.MetricValuesLatest(
		c.GetString("org_id"),
		c.QueryArray("hostId"),
		c.QueryArray("name"),
	)

	doResponse(c, out, err)
}

func (cntr *HostController) SaveMetricValues(c Context) {
	var in []domain.MetricValue
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveMetricValues(
		c.GetString("org_id"),
		in,
	)

	doResponse(c, out, err)
}

func (cntr *HostController) Metadata(c Context) {
	out, err := cntr.Interactor.Metadata(
		c.GetString("org_id"),
		c.Param("hostId"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (cntr *HostController) ListMetadata(c Context) {
	out, err := cntr.Interactor.ListMetadata(
		c.GetString("org_id"),
		c.Param("hostId"),
	)

	doResponse(c, out, err)
}

func (cntr *HostController) SaveMetadata(c Context) {
	var in any
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveMetadata(
		c.GetString("org_id"),
		c.Param("hostId"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (cntr *HostController) DeleteMetadata(c Context) {
	out, err := cntr.Interactor.DeleteMetadata(
		c.GetString("org_id"),
		c.Param("hostId"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}
