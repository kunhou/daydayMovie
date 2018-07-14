package router

import (
	"net/http"

	"github.com/kunhou/TMDB/config"
	"github.com/kunhou/TMDB/provider"

	"github.com/gin-gonic/gin"
	providerHttpDeliver "github.com/kunhou/TMDB/provider/delivery/http"
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

func Setting(pu provider.ProviderUsecase) *gin.Engine {
	router := NewGin()
	router.LoadHTMLGlob(cfg.TempPath)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	pHttpHandler := providerHttpDeliver.NewProviderHttpHandler(pu)
	v1 := router.Group("/v1")
	{
		v1.Any("authentication/:provider/start", providerLogin)
		v1.Any("authentication/:provider/redirect", providerRedirect)
		v1.POST("manual/crawler", pHttpHandler.ManualCrawlerTask)
	}
	return router
}
