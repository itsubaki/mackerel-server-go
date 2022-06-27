package controller

import (
	"net/http"

	"github.com/itsubaki/mackerel-server-go/usecase"
)

func doResponse(c Context, out interface{}, err error) {
	if err == nil {
		c.JSON(http.StatusOK, out)
		return
	}

	switch err.(type) {
	case
		*usecase.ServiceNotFound,
		*usecase.RoleNotFound,
		*usecase.RoleMetadataNotFound,
		*usecase.HostNotFound,
		*usecase.HostMetricNotFound,
		*usecase.HostMetadataNotFound,
		*usecase.ServiceMetricNotFound,
		*usecase.ServiceMetadataNotFound,
		*usecase.AlertNotFound,
		*usecase.InvitationNotFound,
		*usecase.UserNotFound,
		*usecase.ChannelNotFound,
		*usecase.NotificationGroupNotFound:
		c.Status(http.StatusNotFound)
		return
	case
		*usecase.InvalidServiceName,
		*usecase.InvalidRoleName,
		*usecase.InvalidJSONFormat,
		*usecase.HostIsRetired,
		*usecase.MetadataLimitExceeded,
		*usecase.AlertLimitOver:
		c.Status(http.StatusBadRequest)
		return
	case
		*usecase.ServiceMetricPostLimitExceeded:
		c.Status(http.StatusTooManyRequests)
		return
	case
		*usecase.MetadataTooLarge:
		c.Status(http.StatusRequestEntityTooLarge)
	case
		*usecase.PermissionDenied:
		c.Status(http.StatusForbidden)
		return
	default:
		c.Status(http.StatusNotImplemented)
		return
	}
}
