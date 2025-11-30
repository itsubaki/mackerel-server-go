package controller

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/itsubaki/mackerel-server-go/domain"
	"github.com/itsubaki/mackerel-server-go/interface/database"
	"github.com/itsubaki/mackerel-server-go/usecase"
)

type ServiceController struct {
	Interactor *usecase.ServiceInteractor
}

func NewServiceController(handler database.SQLHandler) *ServiceController {
	return &ServiceController{
		Interactor: &usecase.ServiceInteractor{
			NameRule:                regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{1,62}`),
			RoleNameRule:            regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{1,62}`),
			ServiceRepository:       database.NewServiceRepository(handler),
			ServiceMetaRepository:   database.NewServiceMetaRepository(handler),
			ServiceMetricRepository: database.NewServiceMetricRepository(handler),
			RoleRepository:          database.NewRoleRepository(handler),
			RoleMetaRepository:      database.NewRoleMetaRepository(handler),
		},
	}
}

func (cntr *ServiceController) List(c Context) {
	out, err := cntr.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) Save(c Context) {
	var in domain.Service
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

func (cntr *ServiceController) Delete(c Context) {
	out, err := cntr.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) ListRole(c Context) {
	out, err := cntr.Interactor.ListRole(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) SaveRole(c Context) {
	var in domain.Role
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveRole(
		c.GetString("org_id"),
		c.Param("serviceName"),
		&in,
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) DeleteRole(c Context) {
	out, err := cntr.Interactor.DeleteRole(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) MetricNames(c Context) {
	out, err := cntr.Interactor.MetricNames(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) MetricValues(c Context) {
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
		c.Param("serviceName"),
		c.Query("name"),
		from,
		to,
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) SaveMetricValues(c Context) {
	var v []domain.ServiceMetricValue
	if err := c.BindJSON(&v); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveMetricValues(
		c.GetString("org_id"),
		c.Param("serviceName"),
		v,
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) Metadata(c Context) {
	out, err := cntr.Interactor.Metadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) ListMetadata(c Context) {
	out, err := cntr.Interactor.ListMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) SaveMetadata(c Context) {
	var in any
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) DeleteMetadata(c Context) {
	out, err := cntr.Interactor.DeleteMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) RoleMetadata(c Context) {
	out, err := cntr.Interactor.RoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) ListRoleMetadata(c Context) {
	out, err := cntr.Interactor.ListRoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) SaveRoleMetadata(c Context) {
	var in any
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := cntr.Interactor.SaveRoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (cntr *ServiceController) DeleteRoleMetadata(c Context) {
	out, err := cntr.Interactor.DeleteRoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}
