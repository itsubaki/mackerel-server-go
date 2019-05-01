package controllers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type ServiceController struct {
	Interactor *usecase.ServiceInteractor
}

func NewServiceController(sqlHandler database.SQLHandler) *ServiceController {
	return &ServiceController{
		Interactor: &usecase.ServiceInteractor{
			NameRule:     regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
			RoleNameRule: regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
			ServiceRepository: &database.ServiceRepository{
				SQLHandler:          sqlHandler,
				Services:            &domain.Services{},
				ServiceMetadata:     &domain.ServiceMetadataList{},
				ServiceMetricValues: &domain.ServiceMetricValues{},
				Roles:               &domain.Roles{},
				RoleMetadataL:       &domain.RoleMetadataList{},
			},
		},
	}
}

func (s *ServiceController) List(c Context) {
	out, err := s.Interactor.List()
	doResponse(c, out, err)
}

func (s *ServiceController) Save(c Context) {
	var in domain.Service
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(&in)
	doResponse(c, out, err)
}

func (s *ServiceController) Delete(c Context) {
	out, err := s.Interactor.Delete(
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) MetadataList(c Context) {
	out, err := s.Interactor.MetadataList(
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) Metadata(c Context) {
	out, err := s.Interactor.Metadata(
		c.Param("serviceName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) SaveMetadata(c Context) {
	var in interface{}
	if err := c.BindJSON(in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveMetadata(
		c.Param("serviceName"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (s *ServiceController) DeleteMetadata(c Context) {
	out, err := s.Interactor.DeleteMetadata(
		c.Param("serviceName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) RoleList(c Context) {
	out, err := s.Interactor.RoleList(
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

	out, err := s.Interactor.SaveRole(c.Param("serviceName"), &in)
	doResponse(c, out, err)
}

func (s *ServiceController) DeleteRole(c Context) {
	out, err := s.Interactor.DeleteRole(
		c.Param("serviceName"),
		c.Param("roleName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) RoleMetadata(c Context) {
	out, err := s.Interactor.RoleMetadata(
		c.Param("serviveName"),
		c.Param("roleName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) RoleMetadataList(c Context) {
	out, err := s.Interactor.RoleMetadataList(
		c.Param("serviveName"),
		c.Param("roleName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) SaveRoleMetadata(c Context) {
	var in interface{}
	if err := c.BindJSON(in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveRoleMetadata(
		c.Param("serviveName"),
		c.Param("roleName"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (s *ServiceController) DeleteRoleMetadata(c Context) {
	out, err := s.Interactor.DeleteRoleMetadata(
		c.Param("serviveName"),
		c.Param("roleName"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) MetricNames(c Context) {
	out, err := s.Interactor.MetricNames(
		c.Param("serviceName"),
	)

	doResponse(c, out, err)
}

func (s *ServiceController) MetricValues(c Context) {
	from, err := strconv.Atoi(c.Query("from"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	to, err := strconv.Atoi(c.Query("to"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.MetricValues(
		c.Param("serviceName"),
		c.Query("metricName"),
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
		c.Param("serviceName"),
		v,
	)

	doResponse(c, out, err)
}
