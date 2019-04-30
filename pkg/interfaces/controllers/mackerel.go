package controllers

import (
	"regexp"

	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func New() (*Mackerel, error) {
	return &Mackerel{
		ServiceNameRule:         regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
		RoleNameRule:            regexp.MustCompile(`^[a-zA-Z0-9]{1,1}[a-zA-Z0-9_-]{1,62}`),
		MetricNameRule:          regexp.MustCompile(`[a-zA-Z0-9._-]+`),
		ServiceRepository:       database.NewServiceRepository(),
		ServiceRoleRepository:   database.NewServiceRoleRepository(),
		ServiceMetricRepository: database.NewServiceMetricRepository(),
		HostRepository:          database.NewHostRepository(),
		HostMetricRepository:    database.NewHostMetricRepository(),
	}, nil
}

type Mackerel struct {
	ServiceNameRule         *regexp.Regexp
	RoleNameRule            *regexp.Regexp
	MetricNameRule          *regexp.Regexp
	ServiceRepository       *database.ServiceRepository
	ServiceRoleRepository   *database.ServiceRoleRepository
	ServiceMetricRepository *database.ServiceMetricRepository
	HostRepository          *database.HostRepository
	HostMetricRepository    *database.HostMetricRepository
}
