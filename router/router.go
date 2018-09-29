package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kunhou/TMDB/config"
	"github.com/kunhou/TMDB/movie"
	"github.com/kunhou/TMDB/provider"

	movieHttpDeliver "github.com/kunhou/TMDB/movie/delivery/http"
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

func Setting(pu provider.ProviderUsecase, mu movie.MovieUsecase) *gin.Engine {
	router := NewGin()
	router.LoadHTMLGlob(cfg.TempPath)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	pHttpHandler := providerHttpDeliver.NewProviderHttpHandler(pu)
	mHttpHandler := movieHttpDeliver.NewMovieHttpHandler(mu)
	v1 := router.Group("/v1")
	{
		v1.Any("authentication/:provider/start", providerLogin)
		v1.Any("authentication/:provider/redirect", providerRedirect)
		v1.GET("movies", mHttpHandler.MovieList)
		v1.GET("movies/:id", mHttpHandler.MovieDetail)
		v1.GET("tv", mHttpHandler.TVList)
		v1.GET("people", mHttpHandler.PeopleList)
		manual := v1.Group("manual")
		{
			manual.POST("crawler/:type", pHttpHandler.ManualCrawlerTask)
		}
	}
	return router
}
