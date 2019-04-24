package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Must(m Mackerel, err error) Mackerel {
	if err != nil {
		log.Fatalf("new mackerel service: %v", err)
	}
	return m
}

func Handler(m Mackerel) *gin.Engine {
	g := gin.New()

	v0 := g.Group("/api").Group("/v0")

	services := v0.Group("/services")
	services.GET("", func(c *gin.Context) {
		c.JSON(200, m.GetServices())
	})

	return g
}
