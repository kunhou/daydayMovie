package router

import (
	"net/http"

	"github.com/kunhou/TMDB/config"

	"github.com/gin-gonic/gin"
)

var cfg = config.GetConfig()
var secretKey = cfg.JWTSecret

func NewGin() *gin.Engine {
	if cfg.Debug {
		router := gin.Default()
		return router
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	return router
}

func Setting() *gin.Engine {
	router := NewGin()
	router.LoadHTMLGlob("cmd/movie/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	v1 := router.Group("/v1")
	{
		v1.Any("authentication/:provider/start", providerLogin)
		v1.Any("authentication/:provider/redirect", providerRedirect)
	}
	return router
}
