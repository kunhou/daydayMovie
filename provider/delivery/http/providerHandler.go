package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kunhou/TMDB/provider"
)

type HttpProviderHandler struct {
	PUsecase provider.ProviderUsecase
}

func (ph *HttpProviderHandler) ManualCrawlerTask(c *gin.Context) {
	ch := ph.PUsecase.CreateBatchStoreTask()
	go ph.PUsecase.StartCrawler(ch)
	return
}

func NewProviderHttpHandler(pu provider.ProviderUsecase) *HttpProviderHandler {
	handler := &HttpProviderHandler{
		PUsecase: pu,
	}
	return handler
}
