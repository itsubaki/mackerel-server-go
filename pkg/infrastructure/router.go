package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/controllers"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func Default() *gin.Engine {
	return Router(nil)
}

func Router(handler database.SQLHandler) *gin.Engine {
	g := gin.Default()

	auth := controllers.NewAuthController(handler)
	g.Use(func(c *gin.Context) {
		c.Set("Method", c.Request.Method)
		auth.Required(c)
	})

	{
		g.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	}

	v0 := g.Group("/api").Group("/v0")
	{
		services := controllers.NewServiceController(handler)

		s := v0.Group("/services")
		s.GET("", func(c *gin.Context) { services.List(c) })
		s.POST("", func(c *gin.Context) { services.Save(c) })
		s.DELETE("/:serviceName", func(c *gin.Context) { services.Delete(c) })

		s.GET("/:serviceName/roles", func(c *gin.Context) { services.RoleList(c) })
		s.POST("/:serviceName/roles", func(c *gin.Context) { services.SaveRole(c) })
		s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) { services.DeleteRole(c) })

		s.GET("/:serviceName/metric-names", func(c *gin.Context) { services.MetricNames(c) })
		s.GET("/:serviceName/metrics", func(c *gin.Context) { services.MetricValues(c) })
		s.POST("/:serviceName/tsdb", func(c *gin.Context) { services.SaveMetricValues(c) })

		s.GET("/:serviceName/metadata", func(c *gin.Context) { services.MetadataList(c) })
		s.GET("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.Metadata(c) })
		s.PUT("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.SaveMetadata(c) })
		s.DELETE("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.DeleteMetadata(c) })

		s.GET("/:serviceName/roles/:roleName/metadata", func(c *gin.Context) { services.RoleMetadataList(c) })
		s.GET("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.RoleMetadata(c) })
		s.PUT("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.SaveRoleMetadata(c) })
		s.DELETE("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.DeleteRoleMetadata(c) })
	}

	{

		hosts := controllers.NewHostController(handler)

		h := v0.Group("/hosts")
		h.GET("", func(c *gin.Context) { hosts.List(c) })
		h.POST("", func(c *gin.Context) { hosts.Save(c) })

		h.GET("/:hostId", func(c *gin.Context) { hosts.Host(c) })
		h.PUT("/:hostId", func(c *gin.Context) { hosts.Update(c) })
		h.PUT("/:hostId/role-fullnames", func(c *gin.Context) { hosts.RoleFullNames(c) })
		h.POST("/:hostId/status", func(c *gin.Context) { hosts.Status(c) })
		h.POST("/:hostId/retire", func(c *gin.Context) { hosts.Retire(c) })

		h.GET("/:hostId/metric-names", func(c *gin.Context) { hosts.MetricNames(c) })
		h.GET("/:hostId/metrics", func(c *gin.Context) { hosts.MetricValues(c) })
		v0.GET("/tsdb/latest", func(c *gin.Context) { hosts.MetricValuesLatest(c) })
		v0.POST("/tsdb", func(c *gin.Context) { hosts.SaveMetricValues(c) })

		h.GET("/:hostId/metadata", func(c *gin.Context) { hosts.MetadataList(c) })
		h.GET("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.Metadata(c) })
		h.PUT("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.SaveMetadata(c) })
		h.DELETE("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.DeleteMetadata(c) })
	}

	{
		graphs := controllers.NewGraphController(handler)

		d := v0.Group("/graph-defs")
		d.POST("/create", func(c *gin.Context) { graphs.Save(c) })
	}

	{
		reports := controllers.NewCheckReportController(handler)

		r := v0.Group("/monitoring/checks/report")
		r.POST("", func(c *gin.Context) { reports.Save(c) })
	}

	{
		alerts := controllers.NewAlertController(handler)

		a := v0.Group("/alerts")
		a.GET("", func(c *gin.Context) { alerts.List(c) })
		a.POST("/:alertId/close", func(c *gin.Context) { alerts.Close(c) })
	}

	{
		invitations := controllers.NewInvitationController(handler)

		i := v0.Group("/invitations")
		i.GET("", func(c *gin.Context) { invitations.List(c) })
		i.POST("", func(c *gin.Context) { invitations.Save(c) })
		i.POST("/revoke", func(c *gin.Context) { invitations.Revoke(c) })
	}

	{
		users := controllers.NewUserController(handler)

		u := v0.Group("/users")
		u.GET("", func(c *gin.Context) { users.List(c) })
		u.DELETE("/:userId", func(c *gin.Context) { users.Delete(c) })
	}

	{
		org := controllers.NewOrgController(handler)
		o := v0.Group("/org")
		o.GET("", func(c *gin.Context) { org.Org(c) })
	}

	return g
}
