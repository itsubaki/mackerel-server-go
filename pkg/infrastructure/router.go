package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/controllers"
)

func Default() *gin.Engine {

	g := gin.Default()

	g.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v0 := g.Group("/api").Group("/v0")

	handler := NewSQLHandler()
	{
		services := controllers.NewServiceController(handler)

		s := v0.Group("/services")
		s.GET("", func(c *gin.Context) { services.List(c) })
		s.POST("", func(c *gin.Context) { services.Save(c) })
		s.DELETE("/:serviceName", func(c *gin.Context) { services.Delete(c) })

		s.GET("/:serviceName/roles", func(c *gin.Context) { services.RoleList(c) })
		s.POST("/:serviceName/roles", func(c *gin.Context) { services.SaveRole(c) })
		s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) { services.DeleteRole(c) })

		s.GET("/:serviceName/roles/:roleName/metadata", func(c *gin.Context) { services.RoleMetadataList(c) })
		s.GET("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.RoleMetadata(c) })
		s.PUT("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.SaveRoleMetadata(c) })
		s.DELETE("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.DeleteRoleMetadata(c) })

		s.GET("/:serviceName/metric-names", func(c *gin.Context) { services.MetricNames(c) })
		s.GET("/:serviceName/metrics", func(c *gin.Context) { services.MetricValues(c) })
	}

	return g
}

func ApiV0Services(v0 *gin.RouterGroup, m *controllers.Mackerel) {
	s := v0.Group("/services")

	s.POST("/:serviceName/tsdb", func(c *gin.Context) {
		var v ServiceMetricValues
		if err := c.BindJSON(&v); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in := controllers.PostServiceMetricInput{
			ServiceName:        c.Param("serviceName"),
			ServiceMetricValue: v,
		}

		out, err := m.PostServiceMetric(&in)
		doResponse(c, out, err)
	})

	s.GET("/:serviceName/metadata/:namespace", func(c *gin.Context) {
		in := controllers.GetServiceMetadataInput{
			ServiceName: c.Param("serviceName"),
			Namespace:   c.Param("namespace"),
		}

		out, err := m.GetServiceMetadata(&in)
		doResponse(c, out, err)
	})

	s.PUT("/:serviceName/metadata/:namespace", func(c *gin.Context) {
		var v interface{}
		if err := c.BindJSON(&v); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		in := controllers.PutServiceMetadataInput{
			ServiceName: c.Param("serviceName"),
			Namespace:   c.Param("namespace"),
			Metadata:    v,
		}

		out, err := m.PutServiceMetadata(&in)
		doResponse(c, out, err)
	})

	s.DELETE("/:serviceName/metadata/:namespace", func(c *gin.Context) {
		in := controllers.DeleteServiceMetadataInput{
			ServiceName: c.Param("serviceName"),
			Namespace:   c.Param("namespace"),
		}

		out, err := m.DeleteServiceMetadata(&in)
		doResponse(c, out, err)
	})

	s.GET("/:serviceName/metadata", func(c *gin.Context) {
		in := controllers.GetServiceMetadataListInput{
			ServiceName: c.Param("serviceName"),
		}

		out, err := m.GetServiceMetadataList(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Hosts(v0 *gin.RouterGroup, m *controllers.Mackerel) {
	h := v0.Group("/hosts")

	h.POST("", func(c *gin.Context) {
		var in controllers.PostHostInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		out, err := m.PostHost(&in)
		doResponse(c, out, err)
	})

	h.GET("/:hostId", func(c *gin.Context) {
		in := controllers.GetHostInput{
			HostID: c.Param("hostId"),
		}

		out, err := m.GetHost(&in)
		doResponse(c, out, err)
	})

	h.PUT("/:hostId", func(c *gin.Context) {
		var in controllers.PutHostInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.HostID = c.Param("hostId")

		out, err := m.PutHost(&in)
		doResponse(c, out, err)
	})

	h.POST("/:hostId/status", func(c *gin.Context) {
		var in controllers.PostHostStatusInput
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
		var in controllers.PutHostRoleFullNamesInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.HostID = c.Param("hostId")

		out, err := m.PutHostRoleFullNames(&in)
		doResponse(c, out, err)
	})

	h.POST("/:hostId/retire", func(c *gin.Context) {
		var in controllers.PostHostRetiredInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.HostID = c.Param("hostId")

		out, err := m.PostHostRetired(&in)
		doResponse(c, out, err)
	})

	h.GET("", func(c *gin.Context) {
		var in controllers.GetHostsInput
		out, err := m.GetHosts(&in)
		doResponse(c, out, err)
	})

	h.GET("/:hostId/metric-names", func(c *gin.Context) {
		in := controllers.GetHostMetricNamesInput{
			HostID: c.Param("hostId"),
		}

		out, err := m.GetHostMetricNames(&in)
		doResponse(c, out, err)
	})

	h.GET("/:hostId/metrics", func(c *gin.Context) {
		in := controllers.GetHostMetricInput{
			HostID: c.Param("hostId"),
		}

		out, err := m.GetHostMetric(&in)
		doResponse(c, out, err)
	})

	h.GET("/:hostId/metadata/:namespace", func(c *gin.Context) {
		in := controllers.GetHostMetadataInput{
			HostID:    c.Param("hostId"),
			Namespace: c.Param("namespace"),
		}

		out, err := m.GetHostMetadata(&in)
		doResponse(c, out, err)
	})

	h.PUT("/:hostId/metadata/:namespace", func(c *gin.Context) {
		var v interface{}
		if err := c.BindJSON(&v); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		in := controllers.PutHostMetadataInput{
			HostID:    c.Param("hostId"),
			Namespace: c.Param("namespace"),
			Metadata:  v,
		}

		out, err := m.PutHostMetadata(&in)
		doResponse(c, out, err)
	})

	h.DELETE("/:hostId/metadata/:namespace", func(c *gin.Context) {
		in := controllers.DeleteHostMetadataInput{
			HostID:    c.Param("hostId"),
			Namespace: c.Param("namespace"),
		}

		out, err := m.DeleteHostMetadata(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Metrics(v0 *gin.RouterGroup, m *controllers.Mackerel) {
	tsdb := v0.Group("/tsdb")

	tsdb.POST("/", func(c *gin.Context) {
		var v HostMetricValues
		if err := c.BindJSON(&v); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in := controllers.PostHostMetricInput{
			MetricValue: v,
		}

		out, err := m.PostHostMetric(&in)
		doResponse(c, out, err)
	})

	tsdb.GET("/latest", func(c *gin.Context) {
		in := controllers.GetHostMetricLatestInput{}

		out, err := m.GetHostMetricLatest(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Monitoring(v0 *gin.RouterGroup, m *controllers.Mackerel) {
	r := v0.Group("/monitoring/checks/report")

	r.POST("/", func(c *gin.Context) {
		var in controllers.PostCheckReportInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		out, err := m.PostCheckReport(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Alerts(v0 *gin.RouterGroup, m *controllers.Mackerel) {
	r := v0.Group("/alerts")

	r.GET("", func(c *gin.Context) {
		in := controllers.GetAlertInput{}

		out, err := m.GetAlert(&in)
		doResponse(c, out, err)
	})

	r.POST("/:alertId/close", func(c *gin.Context) {
		var in controllers.PostAlertInput
		if err := c.BindJSON(&in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.AlertID = c.Param("alertId")

		out, err := m.PostAlert(&in)
		doResponse(c, out, err)
	})
}

func ApiV0Users(v0 *gin.RouterGroup, m *controllers.Mackerel) {
	r := v0.Group("/users")

	r.GET("", func(c *gin.Context) {
		in := controllers.GetUserInput{}

		out, err := m.GetUser(&in)
		doResponse(c, out, err)
	})

	r.DELETE("/:userId", func(c *gin.Context) {
		in := controllers.DeleteUserInput{
			UserID: c.Param("userId"),
		}

		out, err := m.DeleteUser(&in)
		doResponse(c, out, err)
	})
}
