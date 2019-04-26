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
	g.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			if c.ContentType() != gin.MIMEJSON {
				c.Status(http.StatusBadRequest)
				c.Abort()
			}
		}

		c.Next()
	})

	g.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v0 := g.Group("/api").Group("/v0")
	ApiV0Services(v0, m)
	ApiV0Hosts(v0, m)

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
		in := GetMetricNamesInput{
			ServiceName: c.Param("serviceName"),
		}

		out, err := m.GetMetricNames(&in)
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

	})

	h.POST("/:hostId/status", func(c *gin.Context) {

	})

	h.PUT("/:hostId/role-fullnames", func(c *gin.Context) {

	})

	h.POST("/:hostId/retire", func(c *gin.Context) {

	})

	h.GET("", func(c *gin.Context) {
		out, err := m.GetHosts()
		doResponse(c, out, err)
	})

	h.GET("/:hostId/metric-names", func(c *gin.Context) {

	})
}
