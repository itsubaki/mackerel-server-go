package controllers

import (
	"net/http"

	"github.com/itsubaki/mackerel-api/pkg/usecase"
)

func doResponse(c Context, out interface{}, err error) {
	switch err.(type) {
	case
		*usecase.ServiceNotFound,
		*usecase.RoleNotFound,
		*usecase.HostNotFound,
		*usecase.HostMetricNotFound,
		*usecase.ServiceMetricNotFound,
		*usecase.AlertNotFound,
		*usecase.UserNotFound:
		c.Status(http.StatusNotFound)
		return
	case
		*usecase.InvalidServiceName,
		*usecase.InvalidRoleName,
		*usecase.InvalidJSONFormat,
		*usecase.LimitOver:
		c.Status(http.StatusBadRequest)
		return
	case *usecase.ServiceMetricPostLimitExceeded:
		c.Status(http.StatusTooManyRequests)
		return
	case *usecase.PermissionDenied:
		c.Status(http.StatusForbidden)
		return
	default:
		c.JSON(http.StatusOK, out)
		return
	}
}
