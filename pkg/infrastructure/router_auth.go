package infrastructure

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const XAPIKEY string = "X-Api-Key"

func UseSession(g *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))

	g.Use(sessions.Sessions("GIN_SESSION", store))
	g.Use(func(c *gin.Context) {
		s := sessions.Default(c)
		v := s.Get(XAPIKEY)
		if v != nil {
			c.Request.Header.Set(XAPIKEY, v.(string))
		}

		c.Next()
	})

	// TODO Implement Google Auth
	g.GET("/signin", func(c *gin.Context) {
		s := sessions.Default(c)
		s.Set(XAPIKEY, apikey)
		s.Save()

		c.JSON(http.StatusOK, gin.H{"signin": "ok"})
	})
}
