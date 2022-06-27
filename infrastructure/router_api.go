package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-server-go/interface/controller"
	"github.com/itsubaki/mackerel-server-go/interface/database"
)

func Hosts(v0 *gin.RouterGroup, handler database.SQLHandler) {
	hosts := controller.NewHostController(handler)

	g := v0.Group("/hosts")
	g.GET("", func(c *gin.Context) { hosts.List(c) })
	g.POST("", func(c *gin.Context) { hosts.Save(c) })

	g.GET("/:hostId", func(c *gin.Context) { hosts.Host(c) })
	g.PUT("/:hostId", func(c *gin.Context) { hosts.Update(c) })
	g.PUT("/:hostId/role-fullnames", func(c *gin.Context) { hosts.RoleFullNames(c) })
	g.POST("/:hostId/status", func(c *gin.Context) { hosts.Status(c) })
	g.POST("/:hostId/retire", func(c *gin.Context) { hosts.Retire(c) })

	g.GET("/:hostId/metric-names", func(c *gin.Context) { hosts.MetricNames(c) })
	g.GET("/:hostId/metrics", func(c *gin.Context) { hosts.MetricValues(c) })
	v0.GET("/tsdb/latest", func(c *gin.Context) { hosts.MetricValuesLatest(c) })
	v0.POST("/tsdb", func(c *gin.Context) { hosts.SaveMetricValues(c) })

	g.GET("/:hostId/metadata", func(c *gin.Context) { hosts.ListMetadata(c) })
	g.GET("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.Metadata(c) })
	g.PUT("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.SaveMetadata(c) })
	g.DELETE("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.DeleteMetadata(c) })
}

func Services(v0 *gin.RouterGroup, handler database.SQLHandler) {
	services := controller.NewServiceController(handler)

	g := v0.Group("/services")
	g.GET("", func(c *gin.Context) { services.List(c) })
	g.POST("", func(c *gin.Context) { services.Save(c) })
	g.DELETE("/:serviceName", func(c *gin.Context) { services.Delete(c) })

	g.GET("/:serviceName/roles", func(c *gin.Context) { services.ListRole(c) })
	g.POST("/:serviceName/roles", func(c *gin.Context) { services.SaveRole(c) })
	g.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) { services.DeleteRole(c) })

	g.GET("/:serviceName/metric-names", func(c *gin.Context) { services.MetricNames(c) })
	g.GET("/:serviceName/metrics", func(c *gin.Context) { services.MetricValues(c) })
	g.POST("/:serviceName/tsdb", func(c *gin.Context) { services.SaveMetricValues(c) })

	g.GET("/:serviceName/metadata", func(c *gin.Context) { services.ListMetadata(c) })
	g.GET("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.Metadata(c) })
	g.PUT("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.SaveMetadata(c) })
	g.DELETE("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.DeleteMetadata(c) })

	g.GET("/:serviceName/roles/:roleName/metadata", func(c *gin.Context) { services.ListRoleMetadata(c) })
	g.GET("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.RoleMetadata(c) })
	g.PUT("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.SaveRoleMetadata(c) })
	g.DELETE("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.DeleteRoleMetadata(c) })
}

func Monitors(v0 *gin.RouterGroup, handler database.SQLHandler) {
	monitors := controller.NewMonitorController(handler)

	g := v0.Group("/monitors")
	g.GET("", func(c *gin.Context) { monitors.List(c) })
	g.POST("", func(c *gin.Context) { monitors.Save(c) })
	g.GET("/:monitorId", func(c *gin.Context) { monitors.Monitor(c) })
	g.PUT("/:monitorId", func(c *gin.Context) { monitors.Update(c) })
	g.DELETE("/:monitorId", func(c *gin.Context) { monitors.Delete(c) })
}

func Channels(v0 *gin.RouterGroup, handler database.SQLHandler) {
	channels := controller.NewChannelController(handler)

	g := v0.Group("/channels")
	g.GET("", func(c *gin.Context) { channels.List(c) })
	g.POST("", func(c *gin.Context) { channels.Save(c) })
	g.DELETE("/:channelId", func(c *gin.Context) { channels.Delete(c) })
}

func NotificationGroups(v0 *gin.RouterGroup, handler database.SQLHandler) {
	groups := controller.NewNotificationGroupController(handler)

	g := v0.Group("/notification-groups")
	g.GET("", func(c *gin.Context) { groups.List(c) })
	g.POST("", func(c *gin.Context) { groups.Save(c) })
	g.PUT("/:notificationGroupId", func(c *gin.Context) { groups.Update(c) })
	g.DELETE("/:notificationGroupId", func(c *gin.Context) { groups.Delete(c) })
}

func Graphs(v0 *gin.RouterGroup, handler database.SQLHandler) {
	graphs := controller.NewGraphController(handler)

	g := v0.Group("/graph-defs")
	g.POST("/create", func(c *gin.Context) { graphs.SaveDef(c) })

	a := v0.Group("/graph-annotations")
	a.GET("", func(c *gin.Context) { graphs.List(c) })
	a.POST("", func(c *gin.Context) { graphs.Save(c) })
	a.PUT("/:annotationId", func(c *gin.Context) { graphs.Update(c) })
	a.DELETE("/:annotationId", func(c *gin.Context) { graphs.Delete(c) })
}

func CheckReports(v0 *gin.RouterGroup, handler database.SQLHandler) {
	reports := controller.NewCheckReportController(handler)

	g := v0.Group("/monitoring/checks/report")
	g.POST("", func(c *gin.Context) { reports.Save(c) })
}

func Alerts(v0 *gin.RouterGroup, handler database.SQLHandler) {
	alerts := controller.NewAlertController(handler)

	g := v0.Group("/alerts")
	g.GET("", func(c *gin.Context) { alerts.List(c) })
	g.POST("/:alertId/close", func(c *gin.Context) { alerts.Close(c) })
}

func Dashboards(v0 *gin.RouterGroup, handler database.SQLHandler) {
	dashboards := controller.NewDashboardController(handler)

	g := v0.Group("/dashboards")
	g.GET("", func(c *gin.Context) { dashboards.List(c) })
	g.POST("", func(c *gin.Context) { dashboards.Save(c) })
	g.GET("/:dashboardId", func(c *gin.Context) { dashboards.Dashboard(c) })
	g.PUT("/:dashboardId", func(c *gin.Context) { dashboards.Update(c) })
	g.DELETE("/:dashboardId", func(c *gin.Context) { dashboards.Delete(c) })
}

func Invitations(v0 *gin.RouterGroup, handler database.SQLHandler) {
	invitations := controller.NewInvitationController(handler)

	g := v0.Group("/invitations")
	g.GET("", func(c *gin.Context) { invitations.List(c) })
	g.POST("", func(c *gin.Context) { invitations.Save(c) })
	g.POST("/revoke", func(c *gin.Context) { invitations.Revoke(c) })
}

func Users(v0 *gin.RouterGroup, handler database.SQLHandler) {
	users := controller.NewUserController(handler)

	g := v0.Group("/users")
	g.GET("", func(c *gin.Context) { users.List(c) })
	g.DELETE("/:userId", func(c *gin.Context) { users.Delete(c) })
}

func Organizations(v0 *gin.RouterGroup, handler database.SQLHandler) {
	orgs := controller.NewOrgController(handler)

	g := v0.Group("/org")
	g.GET("", func(c *gin.Context) { orgs.Org(c) })
}

func CheckMonitors(v0 *gin.RouterGroup, handler database.SQLHandler) {
	check := controller.NewCheckMonitorController(handler)

	h := v0.Group("/monitoring/checks/host-metric")
	h.GET("", func(c *gin.Context) { check.HostMetric(c) })

	c := v0.Group("/monitoring/checks/connectivity")
	c.GET("", func(c *gin.Context) { check.Connectivity(c) })

	s := v0.Group("/monitoring/checks/service-metric")
	s.GET("", func(c *gin.Context) { check.ServiceMetric(c) })

	ext := v0.Group("/monitoring/checks/external")
	ext.GET("", func(c *gin.Context) { check.External(c) })

	exp := v0.Group("/monitoring/checks/expression")
	exp.GET("", func(c *gin.Context) { check.Expression(c) })
}

func Downtimes(v0 *gin.RouterGroup, handler database.SQLHandler) {
	downtimes := controller.NewDowntimeController(handler)

	g := v0.Group("/downtimes")
	g.GET("", func(c *gin.Context) { downtimes.List(c) })
	g.POST("", func(c *gin.Context) { downtimes.Save(c) })
	g.GET("/:downtimeId", func(c *gin.Context) { downtimes.Downtime(c) })
	g.PUT("/:downtimeId", func(c *gin.Context) { downtimes.Update(c) })
	g.DELETE("/:downtimeId", func(c *gin.Context) { downtimes.Delete(c) })
}

func UseAPIKey(g *gin.RouterGroup, handler database.SQLHandler) {
	apikeys := controller.NewAPIKeyController(handler)

	g.Use(func(c *gin.Context) {
		key, err := apikeys.APIKey(c)
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
}

func APIv0(g *gin.Engine, handler database.SQLHandler) *gin.Engine {
	v0 := g.Group("/api").Group("/v0")
	UseAPIKey(v0, handler)

	Hosts(v0, handler)
	Services(v0, handler)
	Monitors(v0, handler)
	Channels(v0, handler)
	NotificationGroups(v0, handler)
	Graphs(v0, handler)
	CheckReports(v0, handler)
	Alerts(v0, handler)
	Dashboards(v0, handler)
	Invitations(v0, handler)
	Users(v0, handler)
	Organizations(v0, handler)
	Downtimes(v0, handler)

	// additional
	CheckMonitors(v0, handler)

	return g
}
