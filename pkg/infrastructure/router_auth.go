package infrastructure

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func UseSession(g *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))

	g.Use(sessions.Sessions("GIN_SESSION", store))
	g.Use(func(c *gin.Context) {
		s := sessions.Default(c)
		v := s.Get("X-Api-Key")
		if v != nil {
			c.Request.Header.Set("X-Api-Key", v.(string))
		}

		c.Next()
	})

	// TODO google auth
	g.GET("/signin", func(c *gin.Context) {
		s := sessions.Default(c)
		s.Set("X-Api-Key", apikey)
		s.Save()

		c.JSON(http.StatusOK, gin.H{"signin": "ok"})
	})
}
