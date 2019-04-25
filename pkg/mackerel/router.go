package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Default() *gin.Engine {
	return Router(Must(New()))
}

func Must(m *Mackerel, err error) *Mackerel {
	if err != nil {
		log.Fatalf("new mackerel service: %v", err)
	}
	return m
}

func Router(m *Mackerel) *gin.Engine {
	g := gin.New()

	g.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v0 := g.Group("/api").Group("/v0")
	ApiV0Services(v0, m)
	ApiV0Hosts(v0, m)

	return g
}

func ApiV0Services(v0 *gin.RouterGroup, m *Mackerel) {
	s := v0.Group("/services")

	// https://mackerel.io/api-docs/entry/services#list
	s.GET("", func(c *gin.Context) {
		out, _ := m.GetServices()
		c.JSON(200, out)
	})

	// https://mackerel.io/api-docs/entry/services#create
	s.POST("", func(c *gin.Context) {
		var in PostServiceInput
		if err := parse(c.Request.Body, &in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		out, err := m.PostService(&in)
		switch err.(type) {
		case PermissionDenied:
			c.Status(http.StatusForbidden)
			return
		case InvalidServiceName:
			c.Status(http.StatusBadRequest)
			return
		}

		c.JSON(200, out)
	})

	// https://mackerel.io/api-docs/entry/services#delete
	s.DELETE("/:serviceName", func(c *gin.Context) {
		in := DeleteServiceInput{
			ServiceName: c.Param("serviceName"),
		}

		out, err := m.DeleteService(&in)
		switch err.(type) {
		case PermissionDenied:
			c.Status(http.StatusForbidden)
			return
		case ServiceNotFound:
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(200, out)
	})

	// https://mackerel.io/api-docs/entry/services#rolelist
	s.GET("/:serviceName/roles", func(c *gin.Context) {
		in := GetRolesInput{
			ServiceName: c.Param("serviceName"),
		}

		out, _ := m.GetRoles(&in)
		c.JSON(200, out)
	})

	// https://mackerel.io/api-docs/entry/services#rolecreate
	s.POST("/:serviceName/roles", func(c *gin.Context) {
		var in PostRoleInput
		if err := parse(c.Request.Body, &in); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		in.ServiceName = c.Param("serviceName")

		out, _ := m.PostRole(&in)
		c.JSON(200, out)
	})

	// https://mackerel.io/api-docs/entry/services#roledelete
	s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) {
		in := DeleteRoleInput{
			ServiceName: c.Param("serviceName"),
			RoleName:    c.Param("roleName"),
		}

		out, _ := m.DeleteRole(&in)
		c.JSON(200, out)
	})

	// https://mackerel.io/api-docs/entry/services#metric-names
	s.GET("/:serviceName/metric-names", func(c *gin.Context) {
		in := GetMetricNamesInput{
			ServiceName: c.Param("serviceName"),
		}

		out, _ := m.GetMetricNames(&in)
		c.JSON(200, out)
	})
}

func ApiV0Hosts(v0 *gin.RouterGroup, m *Mackerel) {
	h := v0.Group("/hosts")

	h.GET("", func(c *gin.Context) {
		out, _ := m.GetHosts()
		c.JSON(http.StatusOK, out)
	})

	h.GET("/:hostId", func(c *gin.Context) {
		in := GetHostInput{
			HostID: c.Param("hostId"),
		}

		out, _ := m.GetHost(&in)
		c.JSON(http.StatusOK, out)
	})
}

func parse(body io.ReadCloser, in interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read request body: %v", err)
	}
	defer body.Close()

	if err := json.Unmarshal(b, in); err != nil {
		return fmt.Errorf("unmarshal request body: %v", err)
	}

	return nil
}
