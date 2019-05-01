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

		s.GET("/:serviceName/metadata", func(c *gin.Context) { services.MetadataList(c) })
		s.GET("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.Metadata(c) })
		s.PUT("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.SaveMetadata(c) })
		s.DELETE("/:serviceName/metadata/:namespace", func(c *gin.Context) { services.DeleteMetadata(c) })

		s.GET("/:serviceName/roles", func(c *gin.Context) { services.RoleList(c) })
		s.POST("/:serviceName/roles", func(c *gin.Context) { services.SaveRole(c) })
		s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) { services.DeleteRole(c) })

		s.GET("/:serviceName/roles/:roleName/metadata", func(c *gin.Context) { services.RoleMetadataList(c) })
		s.GET("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.RoleMetadata(c) })
		s.PUT("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.SaveRoleMetadata(c) })
		s.DELETE("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.DeleteRoleMetadata(c) })

		s.GET("/:serviceName/roles/:roleName/metadata", func(c *gin.Context) { services.RoleMetadataList(c) })
		s.GET("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.RoleMetadata(c) })
		s.PUT("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.SaveRoleMetadata(c) })
		s.DELETE("/:serviceName/roles/:roleName/metadata/:namespace", func(c *gin.Context) { services.DeleteRoleMetadata(c) })

		s.GET("/:serviceName/metric-names", func(c *gin.Context) { services.MetricNames(c) })
		s.GET("/:serviceName/metrics", func(c *gin.Context) { services.MetricValues(c) })
		s.POST("/:serviceName/tsdb", func(c *gin.Context) { services.SaveMetricValues(c) })
	}

	return g
}
