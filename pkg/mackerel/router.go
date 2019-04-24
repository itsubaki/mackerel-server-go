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
		c.JSON(200, "ok")
	})

	v0 := g.Group("/api").Group("/v0")

	{
		s := v0.Group("/services")

		// https://mackerel.io/api-docs/entry/services#list
		s.GET("", func(c *gin.Context) {
			out := m.GetServices()
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#create
		s.POST("", func(c *gin.Context) {
			var in PostServiceInput
			if st, err := parse(c.Request.Body, &in); err != nil {
				c.JSON(st, fmt.Errorf("invalid request: %v", err))
				return
			}

			out := m.PostService(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#delete
		s.DELETE("/:serviceName", func(c *gin.Context) {
			in := DeleteServiceInput{
				ServiceName: c.Param("serviceName"),
			}

			out := m.DeleteService(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#rolelist
		s.GET("/:serviceName/roles", func(c *gin.Context) {
			in := GetRolesInput{
				ServiceName: c.Param("serviceName"),
			}

			out := m.GetRoles(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#rolecreate
		s.POST("/:serviceName/roles", func(c *gin.Context) {
			var in PostRoleInput
			if st, err := parse(c.Request.Body, &in); err != nil {
				c.JSON(st, fmt.Errorf("invalid request: %v", err))
				return
			}
			in.ServiceName = c.Param("serviceName")

			out := m.PostRole(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#roledelete
		s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) {
			in := DeleteRoleInput{
				ServiceName: c.Param("serviceName"),
				RoleName:    c.Param("roleName"),
			}

			out := m.DeleteRole(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#metric-names
		s.GET("/:serviceName/metric-names", func(c *gin.Context) {
			in := GetMetricNamesInput{
				ServiceName: c.Param("serviceName"),
			}

			out := m.GetMetricNames(&in)
			c.JSON(out.Status, out)
		})
	}

	{
		h := v0.Group("/hosts")

		h.GET("", func(c *gin.Context) {
			out := m.GetHosts()
			c.JSON(out.Status, out)
		})

		h.GET("/:hostId", func(c *gin.Context) {
			in := GetHostInput{
				HostID: c.Param("hostId"),
			}

			out := m.GetHost(&in)
			c.JSON(out.Status, out)
		})
	}

	return g
}

func parse(body io.ReadCloser, in interface{}) (int, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("read request body: %v", err)
	}
	defer body.Close()

	if err := json.Unmarshal(b, in); err != nil {
		return http.StatusBadRequest, fmt.Errorf("unmarshal request body: %v", err)
	}

	return 200, nil
}
