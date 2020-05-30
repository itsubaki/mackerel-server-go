package controller

import (
	"net/http"
	"strconv"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/itsubaki/mackerel-server-go/pkg/usecase"
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

func (s *HostController) List(c Context) {
	out, err := s.Interactor.List(
		c.GetString("org_id"),
	)

	doResponse(c, out, err)
}

func (s *HostController) Save(c Context) {
	var in domain.Host
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

func (s *HostController) Update(c Context) {
	var in domain.Host
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("hostId")

	out, err := s.Interactor.Save(
		c.GetString("org_id"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *HostController) Host(c Context) {
	out, err := s.Interactor.Host(
		c.GetString("org_id"),
		c.Param("hostId"),
	)

	doResponse(c, out, err)
}

func (s *HostController) Status(c Context) {
	var in domain.HostStatus
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Status(
		c.GetString("org_id"),
		c.Param("hostId"),
		in.Status,
	)

	doResponse(c, out, err)
}

func (s *HostController) RoleFullNames(c Context) {
	var in domain.RoleFullNames
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveRoleFullNames(
		c.GetString("org_id"),
		c.Param("hostId"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *HostController) Retire(c Context) {
	var in domain.HostRetire

	// mkr request dont have empty body.
	//if err := c.BindJSON(&in); err != nil {
	//	c.Status(http.StatusBadRequest)
	//	return
	//}

	out, err := s.Interactor.Retire(
		c.GetString("org_id"),
		c.Param("hostId"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *HostController) MetricNames(c Context) {
	out, err := s.Interactor.MetricNames(
		c.GetString("org_id"),
		c.Param("hostId"),
	)

	doResponse(c, out, err)
}

func (s *HostController) MetricValues(c Context) {
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
		c.Param("hostId"),
		c.Query("name"),
		from,
		to,
	)

	doResponse(c, out, err)
}

func (s *HostController) MetricValuesLatest(c Context) {
	out, err := s.Interactor.MetricValuesLatest(
		c.GetString("org_id"),
		c.QueryArray("hostId"),
		c.QueryArray("name"),
	)

	doResponse(c, out, err)
}

func (s *HostController) SaveMetricValues(c Context) {
	var in []domain.MetricValue
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveMetricValues(
		c.GetString("org_id"),
		in,
	)

	doResponse(c, out, err)
}

func (s *HostController) Metadata(c Context) {
	out, err := s.Interactor.Metadata(
		c.GetString("org_id"),
		c.Param("hostId"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}

func (s *HostController) ListMetadata(c Context) {
	out, err := s.Interactor.ListMetadata(
		c.GetString("org_id"),
		c.Param("hostId"),
	)

	doResponse(c, out, err)
}

func (s *HostController) SaveMetadata(c Context) {
	var in interface{}
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.SaveMetadata(
		c.GetString("org_id"),
		c.Param("hostId"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (s *HostController) DeleteMetadata(c Context) {
	out, err := s.Interactor.DeleteMetadata(
		c.GetString("org_id"),
		c.Param("hostId"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}
