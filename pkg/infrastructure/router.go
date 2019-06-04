package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/controllers"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func Router(handler database.SQLHandler) *gin.Engine {
	g := gin.Default()

	auth := controllers.NewAuthController(handler)
	g.Use(func(c *gin.Context) {
		key, err := auth.XAPIKey(c)
		if err != nil {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}

		if c.Request.Method != http.MethodGet && !key.Write {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}

		c.Set("org_id", key.OrgID)
		c.Next()
	})

	{
		g.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	}

	v0 := g.Group("/api").Group("/v0")
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
		monitors := controllers.NewMonitorController(handler)

		m := v0.Group("/monitors")
		m.GET("", func(c *gin.Context) { monitors.List(c) })
		m.POST("", func(c *gin.Context) { monitors.Save(c) })
		m.GET("/:monitorId", func(c *gin.Context) { monitors.Monitor(c) })
		m.PUT("/:monitorId", func(c *gin.Context) { monitors.Update(c) })
		m.DELETE("/:monitorId", func(c *gin.Context) { monitors.Delete(c) })
	}

	{
		channels := controllers.NewChannelController(handler)

		c := v0.Group("/channels")
		c.GET("", func(c *gin.Context) { channels.List(c) })
		c.POST("", func(c *gin.Context) { channels.Save(c) })
		c.DELETE("/:channelId", func(c *gin.Context) { channels.Delete(c) })
	}

	{
		groups := controllers.NewNotificationGroupController(handler)

		g := v0.Group("/notification-groups")
		g.GET("", func(c *gin.Context) { groups.List(c) })
		g.POST("", func(c *gin.Context) { groups.Save(c) })
		g.PUT("/:notificationGroupId", func(c *gin.Context) { groups.Update(c) })
		g.DELETE("/:notificationGroupId", func(c *gin.Context) { groups.Delete(c) })
	}

	{
		graphs := controllers.NewGraphController(handler)

		d := v0.Group("/graph-defs")
		d.POST("/create", func(c *gin.Context) { graphs.SaveDef(c) })

		a := v0.Group("/graph-annotations")
		a.GET("", func(c *gin.Context) {})
		a.POST("", func(c *gin.Context) {})
		a.PUT("/:annotationId", func(c *gin.Context) {})
		a.DELETE("/:annotationId", func(c *gin.Context) {})
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
		dashboard := controllers.NewDashboardController(handler)

		d := v0.Group("/dashboards")
		d.GET("", func(c *gin.Context) { dashboard.List(c) })
		d.POST("", func(c *gin.Context) { dashboard.Save(c) })
		d.GET("/:dashboardId", func(c *gin.Context) { dashboard.Dashboard(c) })
		d.PUT("/:dashboardId", func(c *gin.Context) { dashboard.Update(c) })
		d.DELETE("/:dashboardId", func(c *gin.Context) { dashboard.Delete(c) })
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
