package controllers

import (
	"net/http"
	"strconv"

	"github.com/itsubaki/mackerel-api/pkg/domain"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/memory"
	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

type HostController struct {
	Interactor *usecase.HostInteractor
}

func NewHostController(handler database.SQLHandler) *HostController {
	var repo usecase.HostRepository
	repo = memory.NewHostRepository()
	if handler != nil {
		repo = database.NewHostRepository(handler)
	}

	return &HostController{
		Interactor: &usecase.HostInteractor{
			HostRepository: repo,
		},
	}
}

func (s *HostController) List(c Context) {
	out, err := s.Interactor.List()
	doResponse(c, out, err)
}

func (s *HostController) Save(c Context) {
	var in domain.Host
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Save(&in)
	doResponse(c, out, err)
}

func (s *HostController) Update(c Context) {
	var in domain.Host
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	in.ID = c.Param("hostId")

	out, err := s.Interactor.Update(&in)
	doResponse(c, out, err)
}

func (s *HostController) Host(c Context) {
	out, err := s.Interactor.Host(
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
		c.Param("hostId"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *HostController) Retire(c Context) {
	var in domain.HostRetire
	if err := c.BindJSON(&in); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	out, err := s.Interactor.Retire(
		c.Param("hostId"),
		&in,
	)

	doResponse(c, out, err)
}

func (s *HostController) MetricNames(c Context) {
	out, err := s.Interactor.MetricNames(
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
		c.Param("hostId"),
		c.Query("name"),
		from,
		to,
	)

	doResponse(c, out, err)
}

func (s *HostController) MetricValuesLatest(c Context) {
	out, err := s.Interactor.MetricValuesLatest(
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

	out, err := s.Interactor.SaveMetricValues(in)
	doResponse(c, out, err)
}

func (s *HostController) MetadataList(c Context) {
	out, err := s.Interactor.MetadataList(
		c.Param("hostId"),
	)

	doResponse(c, out, err)
}

func (s *HostController) Metadata(c Context) {
	out, err := s.Interactor.Metadata(
		c.Param("hostId"),
		c.Param("namespace"),
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
		c.Param("hostId"),
		c.Param("namespace"),
		in,
	)

	doResponse(c, out, err)
}

func (s *HostController) DeleteMetadata(c Context) {
	out, err := s.Interactor.DeleteMetadata(
		c.Param("hostId"),
		c.Param("namespace"),
	)

	doResponse(c, out, err)
}
