package controllers

import (
	"net/http"
	"regexp"

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
			ServiceNameRule:     regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
			ServiceRoleNameRule: regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
			ServiceRepository: &database.ServiceRepository{
				SQLHandler:          sqlHandler,
				Services:            domain.Services{},
				ServiceMetadata:     domain.ServiceMetadataList{},
				ServiceMetricValues: domain.ServiceMetricValues{},
				Roles:               domain.Roles{},
				RoleMetadata:        domain.RoleMetadataList{},
			},
		},
	}
}

func (s *ServiceController) MetricNames(c Context) {
	out, err := s.Interactor.MetricNames(c.Param("serviceName"))
	doResponse(c, out, err)
}

func (s *ServiceController) MetricValues(c Context) {
	out, err := s.Interactor.MetricValues(c.Param("serviceName"), c.Param("metricName"), 0, 0)
	doResponse(c, out, err)
}

func (s *ServiceController) RoleList(c Context) {
	out, err := s.Interactor.RoleList(c.Param("serviceName"))
	doResponse(c, out, err)
}

func (s *ServiceController) SaveRole(c Context) {
	var in domain.Role
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveRole(&in)
	doResponse(c, out, err)
}

func (s *ServiceController) DeleteRole(c Context) {
	out, err := s.Interactor.DeleteRole(c.Param("serviceName"), c.Param("roleName"))
	doResponse(c, out, err)
}

func (s *ServiceController) RoleMetadata(c Context) {
	out, err := s.Interactor.RoleMetadata(c.Param("serviveName"), c.Param("roleName"), c.Param("namespace"))
	doResponse(c, out, err)
}

func (s *ServiceController) RoleMetadataList(c Context) {
	type Response struct {
		Metadata domain.RoleMetadataList `json:"metadata"`
	}

	out, err := s.Interactor.RoleMetadataList(c.Param("serviveName"), c.Param("roleName"))
	doResponse(c, &Response{Metadata: out}, err)
}

func (s *ServiceController) SaveRoleMetadata(c Context) {
	var in interface{}
	if err := c.BindJSON(in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	type Response struct {
		Success bool `json:"success"`
	}

	if err := s.Interactor.SaveRoleMetadata(c.Param("serviveName"), c.Param("roleName"), c.Param("namespace"), in); err != nil {
		doResponse(c, &Response{Success: false}, err)
		return
	}

	doResponse(c, &Response{Success: true}, nil)
}

func (s *ServiceController) DeleteRoleMetadata(c Context) {
	type Response struct {
		Success bool `json:"success"`
	}

	if err := s.Interactor.DeleteRoleMetadata(c.Param("serviveName"), c.Param("roleName"), c.Param("namespace")); err != nil {
		doResponse(c, &Response{Success: false}, err)
		return
	}

	doResponse(c, &Response{Success: true}, nil)
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

func (s *ServiceController) List(c Context) {
	type Response struct {
		Services domain.Services `json:"services"`
	}

	out, err := s.Interactor.List()
	doResponse(c, &Response{Services: out}, err)
}

func (s *ServiceController) Delete(c Context) {
	out, err := s.Interactor.Delete(c.Param("serviceName"))
	doResponse(c, out, err)
}
