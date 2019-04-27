package mackerel

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Default() *gin.Engine {
	return Router(gin.Default(), Must(New()))
}

func Must(m *Mackerel, err error) *Mackerel {
	if err != nil {
		log.Fatalf("new mackerel service: %v", err)
	}
	return m
}

func Router(g *gin.Engine, m *Mackerel) *gin.Engine {
	g.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v0 := g.Group("/api").Group("/v0")
	ApiV0Services(v0, m)
	ApiV0Hosts(v0, m)
	ApiV0Metrics(v0, m)
	ApiV0Monitoring(v0, m)

	return g
}

func doResponse(c *gin.Context, out interface{}, err error) {
	switch err.(type) {
	case *ServiceNotFound:
		c.Status(http.StatusNotFound)
		return
	case *InvalidServiceName:
		c.Status(http.StatusBadRequest)
		return
	case *RoleNotFound:
		c.Status(http.StatusNotFound)
		return
	case *InvalidRoleName:
		c.Status(http.StatusBadRequest)
		return
	case *HostNotFound:
		c.Status(http.StatusNotFound)
		return
	case *InvalidJSONFormat:
		c.Status(http.StatusBadRequest)
		return
	case *HostMetricNotFound:
		c.Status(http.StatusNotFound)
		return
	case *ServiceMetricNotFound:
		c.Status(http.StatusNotFound)
		return
	case *ServiceMetricPostLimitExceeded:
		c.Status(http.StatusTooManyRequests)
		return
	case *PermissionDenied:
		c.Status(http.StatusForbidden)
		return
	default:
		c.JSON(http.StatusOK, out)
		return
	}
}

func ApiV0Services(v0 *gin.RouterGroup, m *Mackerel) {
	s := v0.Group("/services")

	// https://mackerel.io/api-docs/entry/services#list
	s.GET("", func(c *gin.Context) {
		out, err := m.GetServices(&GetServicesInput{})
		doResponse(c, out, err)
	})

	// https://mackerel.io/api-docs/entry/services#create
	s.POST("", func(c *gin.Context) {
		var in PostServiceInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		out, err := m.PostService(&in)
		doResponse(c, out, err)
	})

	// https://mackerel.io/api-docs/entry/services#delete
	s.DELETE("/:serviceName", func(c *gin.Context) {
		in := DeleteServiceInput{
			ServiceName: c.Param("serviceName"),
		}

		out, err := m.DeleteService(&in)
		doResponse(c, out, err)
	})

	// https://mackerel.io/api-docs/entry/services#rolelist
	s.GET("/:serviceName/roles", func(c *gin.Context) {
		in := GetRolesInput{
			ServiceName: c.Param("serviceName"),
		}

		out, err := m.GetRoles(&in)
		doResponse(c, out, err)
	})

	// https://mackerel.io/api-docs/entry/services#rolecreate
	s.POST("/:serviceName/roles", func(c *gin.Context) {
		var in PostRoleInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.ServiceName = c.Param("serviceName")

		out, err := m.PostRole(&in)
		doResponse(c, out, err)
	})

	// https://mackerel.io/api-docs/entry/services#roledelete
	s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) {
		in := DeleteRoleInput{
			ServiceName: c.Param("serviceName"),
			RoleName:    c.Param("roleName"),
		}

		out, err := m.DeleteRole(&in)
		doResponse(c, out, err)
	})

	// https://mackerel.io/api-docs/entry/services#metric-names
	s.GET("/:serviceName/metric-names", func(c *gin.Context) {
		in := GetServiceMetricNamesInput{
			ServiceName: c.Param("serviceName"),
		}

		out, err := m.GetServiceMetricNames(&in)
		doResponse(c, out, err)
	})

	s.POST("/:serviceName/tsdb", func(c *gin.Context) {
		var v []ServiceMetricValue
		if err := c.BindJSON(&v); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in := PostServiceMetricInput{
			ServiceName:        c.Param("serviceName"),
			ServiceMetricValue: v,
		}

		out, err := m.PostServiceMetric(&in)
		doResponse(c, out, err)
	})

	s.GET("/:serviceName/metrics", func(c *gin.Context) {
		in := GetServiceMetricInput{
			ServiceName: c.Param("serviceName"),
		}

		out, err := m.GetServiceMetric(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Hosts(v0 *gin.RouterGroup, m *Mackerel) {
	h := v0.Group("/hosts")

	h.POST("", func(c *gin.Context) {
		var in PostHostInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		out, err := m.PostHost(&in)
		doResponse(c, out, err)
	})

	h.GET("/:hostId", func(c *gin.Context) {
		in := GetHostInput{
			HostID: c.Param("hostId"),
		}

		out, err := m.GetHost(&in)
		doResponse(c, out, err)
	})

	h.PUT("/:hostId", func(c *gin.Context) {
		var in PutHostInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.HostID = c.Param("hostId")

		out, err := m.PutHost(&in)
		doResponse(c, out, err)
	})

	h.POST("/:hostId/status", func(c *gin.Context) {
		var in PostHostStatusInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.HostID = c.Param("hostId")

		if in.Status != "standby" &&
			in.Status != "working" &&
			in.Status != "maintenance" &&
			in.Status != "poweroff" {
			c.Status(http.StatusBadRequest)
			return
		}

		out, err := m.PostHostStatus(&in)
		doResponse(c, out, err)
	})

	h.PUT("/:hostId/role-fullnames", func(c *gin.Context) {
		var in PutHostRoleFullNamesInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.HostID = c.Param("hostId")

		out, err := m.PutHostRoleFullNames(&in)
		doResponse(c, out, err)
	})

	h.POST("/:hostId/retire", func(c *gin.Context) {
		var in PostHostRetiredInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.HostID = c.Param("hostId")

		out, err := m.PostHostRetired(&in)
		doResponse(c, out, err)
	})

	h.GET("", func(c *gin.Context) {
		var in GetHostsInput
		out, err := m.GetHosts(&in)
		doResponse(c, out, err)
	})

	h.GET("/:hostId/metric-names", func(c *gin.Context) {
		in := GetHostMetricNamesInput{
			HostID: c.Param("hostId"),
		}

		out, err := m.GetHostMetricNames(&in)
		doResponse(c, out, err)
	})

	h.GET("/:hostId/metrics", func(c *gin.Context) {
		in := GetHostMetricInput{
			HostID: c.Param("hostId"),
		}

		out, err := m.GetHostMetric(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Metrics(v0 *gin.RouterGroup, m *Mackerel) {
	tsdb := v0.Group("/tsdb")

	tsdb.POST("/", func(c *gin.Context) {
		var v []HostMetricValue
		if err := c.BindJSON(&v); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in := PostHostMetricInput{
			MetricValue: v,
		}

		out, err := m.PostHostMetric(&in)
		doResponse(c, out, err)
	})

	tsdb.GET("/latest", func(c *gin.Context) {
		in := GetHostMetricLatestInput{}

		out, err := m.GetHostMetricLatest(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Monitoring(v0 *gin.RouterGroup, m *Mackerel) {
	r := v0.Group("/monitoring/checks/report")

	r.POST("/", func(c *gin.Context) {
		var in PostCheckReportInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		out, err := m.PostCheckReport(&in)
		doResponse(c, out, err)
	})
}
