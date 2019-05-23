package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/controllers"
	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
)

func Default() *gin.Engine {
	return Router(nil)
}

// mysql> explain select * from xapikey where xwrite='1' and x_api_key='2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb';
// +----+-------------+---------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// | id | select_type | table   | partitions | type  | possible_keys | key     | key_len | ref   | rows | filtered | Extra |
// +----+-------------+---------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// |  1 | SIMPLE      | xapikey | NULL       | const | PRIMARY       | PRIMARY | 182     | const |    1 |   100.00 | NULL  |
// +----+-------------+---------+------------+-------+---------------+---------+---------+-------+------+----------+-------+
// 1 row in set, 1 warning (0.00 sec)
func auth(handler database.SQLHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var org string
		write := false
		if err := handler.Transact(func(tx database.Tx) error {
			row := tx.QueryRow(
				`
				select org, xwrite from xapikey where x_api_key=?
				`,
				c.GetHeader("X-Api-Key"),
			)

			if err := row.Scan(
				&org,
				&write,
			); err != nil {
				return fmt.Errorf("select * from xapikey: %v", err)
			}

			return nil
		}); err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
		}

		if c.Request.Method != http.MethodGet && !write {
			c.Status(http.StatusForbidden)
			c.Abort()
		}

		c.Next()
	}
}

func Router(handler database.SQLHandler) *gin.Engine {
	g := gin.Default()
	g.Use(auth(handler))

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
