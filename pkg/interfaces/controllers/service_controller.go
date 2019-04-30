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
	out, err := s.Interactor.List()
	doResponse(c, out, err)
}

func (s *ServiceController) Delete(c Context) {
	out, err := s.Interactor.Delete(c.Param("serviceName"))
	doResponse(c, out, err)
}
