package controller

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
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

func (s *ServiceController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) Save(c Context) {
	var in domain.Service
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

func (s *ServiceController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) ListRole(c Context) {
	out, err := s.Interactor.ListRole(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) SaveRole(c Context) {
	var in domain.Role
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveRole(
		c.GetString("org_id"),
		c.Param("serviceName"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *ServiceController) DeleteRole(c Context) {
	out, err := s.Interactor.DeleteRole(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) MetricNames(c Context) {
	out, err := s.Interactor.MetricNames(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) MetricValues(c Context) {
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

	out, err := s.Interactor.MetricValues(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Query("name"),
		from,
		to,
	)

	doResponse(c, out, err)
}

func (s *ServiceController) SaveMetricValues(c Context) {
	var v []domain.ServiceMetricValue
	if err := c.BindJSON(&v); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveMetricValues(
		c.GetString("org_id"),
		c.Param("serviceName"),
		v,
	)

	doResponse(c, out, err)
}

func (s *ServiceController) Metadata(c Context) {
	out, err := s.Interactor.Metadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) ListMetadata(c Context) {
	out, err := s.Interactor.ListMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) SaveMetadata(c Context) {
	var in interface{}
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (s *ServiceController) DeleteMetadata(c Context) {
	out, err := s.Interactor.DeleteMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) RoleMetadata(c Context) {
	out, err := s.Interactor.RoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) ListRoleMetadata(c Context) {
	out, err := s.Interactor.ListRoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) SaveRoleMetadata(c Context) {
	var in interface{}
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveRoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (s *ServiceController) DeleteRoleMetadata(c Context) {
	out, err := s.Interactor.DeleteRoleMetadata(
		c.GetString("org_id"),
		c.Param("serviceName"),
		c.Param("roleName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}
