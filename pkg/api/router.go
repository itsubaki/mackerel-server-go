package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/internal/services"
)

func Must(m Mackerel, err error) Mackerel {
	if err != nil {
		log.Fatalf("new mackerel service: %v", err)
	}
	return m
}

func Router(m Mackerel) *gin.Engine {
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
			var in services.PostServicesInput
			if st, err := parse(c.Request.Body, &in); err != nil {
				c.JSON(st, fmt.Errorf("invalid request: %v", err))
				return
			}

			out := m.PostServices(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#delete
		s.DELETE("/:serviceName", func(c *gin.Context) {
			in := services.DeleteServicesInput{
				ServiceName: c.Param("serviceName"),
			}

			out := m.DeleteServices(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#rolelist
		s.GET("/:serviceName/roles", func(c *gin.Context) {
			in := services.GetRolesInput{
				ServiceName: c.Param("serviceName"),
			}

			out := m.GetRoles(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#rolecreate
		s.POST("/:serviceName/roles", func(c *gin.Context) {
			var in services.PostRolesInput
			if st, err := parse(c.Request.Body, &in); err != nil {
				c.JSON(st, fmt.Errorf("invalid request: %v", err))
				return
			}
			in.ServiceName = c.Param("serviceName")

			out := m.PostRoles(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#roledelete
		s.DELETE("/:serviceName/roles/:roleName", func(c *gin.Context) {
			in := services.DeleteRolesInput{
				ServiceName: c.Param("serviceName"),
				RoleName:    c.Param("roleName"),
			}

			out := m.DeleteRoles(&in)
			c.JSON(out.Status, out)
		})

		// https://mackerel.io/api-docs/entry/services#metric-names
		s.GET("/:serviceName/metric-names", func(c *gin.Context) {
			in := services.GetMetricNamesInput{
				ServiceName: c.Param("serviceName"),
			}

			out := m.GetMetricNames(&in)
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
