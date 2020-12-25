package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/controller"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
)

func Hosts(v0 *gin.RouterGroup, handler database.SQLHandler) {
	hosts := controller.NewHostController(handler)

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

	h.GET("/:hostId/metadata", func(c *gin.Context) { hosts.ListMetadata(c) })
	h.GET("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.Metadata(c) })
	h.PUT("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.SaveMetadata(c) })
	h.DELETE("/:hostId/metadata/:namespace", func(c *gin.Context) { hosts.DeleteMetadata(c) })
}

func Services(v0 *gin.RouterGroup, handler database.SQLHandler) {
	services := controller.NewServiceController(handler)

	s := v0.Group("/services")
	s.GET("", func(c *gin.Context) { services.List(c) })
	s.POST("", func(c *gin.Context) { services.Save(c) })
	s.DELETE("/:serviceName", func(c *gin.Context) { services.Delete(c) })

	s.GET("/:serviceName/roles", func(c *gin.Context) { services.ListRole(c) })
	s.POST("/:serviceName/roles", func(c *gin.Context) { services.SaveRole(c) })
	s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) { services.DeleteRole(c) })

	s.GET("/:serviceName/metric-names", func(c *gin.Context) { services.MetricNames(c) })
	s.GET("/:serviceName/metrics", func(c *gin.Context) { services.MetricValues(c) })
	s.POST("/:serviceName/tsdb", func(c *gin.Context) { services.SaveMetricValues(c) })

	s.GET("/:serviceName/metadata", func(c *gin.Context) { services.ListMetadata(c) })
	s.GET("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.Metadata(c) })
	s.PUT("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.SaveMetadata(c) })
	s.DELETE("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.DeleteMetadata(c) })

	s.GET("/:serviceName/roles/:roleName/metadata", func(c *gin.Context) { services.ListRoleMetadata(c) })
	s.GET("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.RoleMetadata(c) })
	s.PUT("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.SaveRoleMetadata(c) })
	s.DELETE("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.DeleteRoleMetadata(c) })
}

func Monitors(v0 *gin.RouterGroup, handler database.SQLHandler) {
	monitors := controller.NewMonitorController(handler)

	m := v0.Group("/monitors")
	m.GET("", func(c *gin.Context) { monitors.List(c) })
	m.POST("", func(c *gin.Context) { monitors.Save(c) })
	m.GET("/:monitorId", func(c *gin.Context) { monitors.Monitor(c) })
	m.PUT("/:monitorId", func(c *gin.Context) { monitors.Update(c) })
	m.DELETE("/:monitorId", func(c *gin.Context) { monitors.Delete(c) })
}

func Channels(v0 *gin.RouterGroup, handler database.SQLHandler) {
	channels := controller.NewChannelController(handler)

	c := v0.Group("/channels")
	c.GET("", func(c *gin.Context) { channels.List(c) })
	c.POST("", func(c *gin.Context) { channels.Save(c) })
	c.DELETE("/:channelId", func(c *gin.Context) { channels.Delete(c) })
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

	d := v0.Group("/graph-defs")
	d.POST("/create", func(c *gin.Context) { graphs.SaveDef(c) })

	a := v0.Group("/graph-annotations")
	a.GET("", func(c *gin.Context) { graphs.List(c) })
	a.POST("", func(c *gin.Context) { graphs.Save(c) })
	a.PUT("/:annotationId", func(c *gin.Context) { graphs.Update(c) })
	a.DELETE("/:annotationId", func(c *gin.Context) { graphs.Delete(c) })
}

func CheckReports(v0 *gin.RouterGroup, handler database.SQLHandler) {
	reports := controller.NewCheckReportController(handler)

	r := v0.Group("/monitoring/checks/report")
	r.POST("", func(c *gin.Context) { reports.Save(c) })
}

func Alerts(v0 *gin.RouterGroup, handler database.SQLHandler) {
	alerts := controller.NewAlertController(handler)

	a := v0.Group("/alerts")
	a.GET("", func(c *gin.Context) { alerts.List(c) })
	a.POST("/:alertId/close", func(c *gin.Context) { alerts.Close(c) })
}

func Dashboards(v0 *gin.RouterGroup, handler database.SQLHandler) {
	dashboard := controller.NewDashboardController(handler)

	d := v0.Group("/dashboards")
	d.GET("", func(c *gin.Context) { dashboard.List(c) })
	d.POST("", func(c *gin.Context) { dashboard.Save(c) })
	d.GET("/:dashboardId", func(c *gin.Context) { dashboard.Dashboard(c) })
	d.PUT("/:dashboardId", func(c *gin.Context) { dashboard.Update(c) })
	d.DELETE("/:dashboardId", func(c *gin.Context) { dashboard.Delete(c) })
}

func Invitations(v0 *gin.RouterGroup, handler database.SQLHandler) {
	invitations := controller.NewInvitationController(handler)

	i := v0.Group("/invitations")
	i.GET("", func(c *gin.Context) { invitations.List(c) })
	i.POST("", func(c *gin.Context) { invitations.Save(c) })
	i.POST("/revoke", func(c *gin.Context) { invitations.Revoke(c) })
}

func Users(v0 *gin.RouterGroup, handler database.SQLHandler) {
	users := controller.NewUserController(handler)

	u := v0.Group("/users")
	u.GET("", func(c *gin.Context) { users.List(c) })
	u.DELETE("/:userId", func(c *gin.Context) { users.Delete(c) })
}

func Organizations(v0 *gin.RouterGroup, handler database.SQLHandler) {
	org := controller.NewOrgController(handler)

	o := v0.Group("/org")
	o.GET("", func(c *gin.Context) { org.Org(c) })
}

func CheckMonitors(v0 *gin.RouterGroup, handler database.SQLHandler) {
	check := controller.NewCheckMonitorController(handler)

	host := v0.Group("/monitoring/checks/host-metric")
	host.GET("", func(c *gin.Context) { check.HostMetric(c) })

	conn := v0.Group("/monitoring/checks/connectivity")
	conn.GET("", func(c *gin.Context) { check.Connectivity(c) })

	service := v0.Group("/monitoring/checks/service-metric")
	service.GET("", func(c *gin.Context) { check.ServiceMetric(c) })

	ext := v0.Group("/monitoring/checks/external")
	ext.GET("", func(c *gin.Context) { check.External(c) })

	exp := v0.Group("/monitoring/checks/expression")
	exp.GET("", func(c *gin.Context) { check.Expression(c) })
}

func Downtimes(v0 *gin.RouterGroup, handler database.SQLHandler) {
	downtime := controller.NewDowntimeController(handler)

	d := v0.Group("/downtimes")
	d.GET("", func(c *gin.Context) { downtime.List(c) })
	d.POST("", func(c *gin.Context) { downtime.Save(c) })
	d.GET("/:downtimeId", func(c *gin.Context) { downtime.Downtime(c) })
	d.PUT("/:downtimeId", func(c *gin.Context) { downtime.Update(c) })
	d.DELETE("/:downtimeId", func(c *gin.Context) { downtime.Delete(c) })
}

func UseAPIKey(g *gin.RouterGroup, handler database.SQLHandler) {
	apikey := controller.NewAPIKeyController(handler)

	g.Use(func(c *gin.Context) {
		key, err := apikey.APIKey(c)
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
	// Monitors(v0, handler)
	// Channels(v0, handler)
	// NotificationGroups(v0, handler)
	Graphs(v0, handler)
	CheckReports(v0, handler)
	Alerts(v0, handler)
	Dashboards(v0, handler)
	Invitations(v0, handler)
	Users(v0, handler)
	Organizations(v0, handler)
	Downtimes(v0, handler)

	// additional
	// CheckMonitors(v0, handler)

	return g
}
