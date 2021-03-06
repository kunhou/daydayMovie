package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/provider"
)

type HttpProviderHandler struct {
	PUsecase provider.ProviderUsecase
}

func (ph *HttpProviderHandler) ManualCrawlerTask(c *gin.Context) {
	crawlerType := c.Param("type")
	if strings.EqualFold(crawlerType, "movie") {
		log.Info("Manual crawler movie")
		ch := ph.PUsecase.CreateBatchStoreMovieTask()
		go ph.PUsecase.StartCrawlerMovie(ch)
	} else if strings.EqualFold(crawlerType, "person") {
		log.Info("Manual crawler person")
		ch := ph.PUsecase.CreateBatchStorePersonTask()
		go ph.PUsecase.StartCrawlerPerson(ch)
	} else if strings.EqualFold(crawlerType, "tv") {
		log.Info("Manual crawler tv")
		ch := ph.PUsecase.CreateStoreTVTask()
		go ph.PUsecase.StartCrawlerTV(ch)
	} else if strings.EqualFold(crawlerType, "credit") {
		log.Info("Manual crawler credit")
		ch := ph.PUsecase.CreateStoreCreditTask()
		pch := ph.PUsecase.CreateBatchStorePersonTask()
		go ph.PUsecase.StartCrawlerCredit(ch, pch)
	} else if strings.EqualFold(crawlerType, "popular_movie") {
		log.Info("Manual crawler popular movie")
		ch := ph.PUsecase.CreateStoreCreditTask()
		pch := ph.PUsecase.CreateBatchStorePersonTask()
		mch := ph.PUsecase.CreateBatchStoreMovieTask()
		go ph.PUsecase.StartCrawlerPopularMovie(ch, pch, mch)
	} else if strings.EqualFold(crawlerType, "popular_tv") {
		log.Info("Manual crawler popular movie")
		ch := ph.PUsecase.CreateStoreCreditTask()
		pch := ph.PUsecase.CreateBatchStorePersonTask()
		tvch := ph.PUsecase.CreateStoreTVTask()
		go ph.PUsecase.StartCrawlerPopularTV(ch, pch, tvch)
	}

	return
}

func NewProviderHttpHandler(pu provider.ProviderUsecase) *HttpProviderHandler {
	handler := &HttpProviderHandler{
		PUsecase: pu,
	}
	return handler
}
