package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/controllers"
)

func doResponse(c *gin.Context, out interface{}, err error) {
	switch err.(type) {
	case
		*controllers.ServiceNotFound,
		*controllers.RoleNotFound,
		*controllers.HostNotFound,
		*controllers.HostMetricNotFound,
		*controllers.ServiceMetricNotFound,
		*controllers.UserNotFound:
		c.Status(http.StatusNotFound)
		return
	case
		*controllers.InvalidServiceName,
		*controllers.InvalidRoleName,
		*controllers.InvalidJSONFormat:
		c.Status(http.StatusBadRequest)
		return
	case *controllers.ServiceMetricPostLimitExceeded:
		c.Status(http.StatusTooManyRequests)
		return
	case *controllers.PermissionDenied:
		c.Status(http.StatusForbidden)
		return
	default:
		c.JSON(http.StatusOK, out)
		return
	}
}
