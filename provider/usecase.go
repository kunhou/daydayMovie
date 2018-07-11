package provider

import "github.com/kunhou/TMDB/models"

type ProviderUsecase interface {
	StartCrawler(ch chan *models.Movie)
	CreateBatchStoreTask() chan *models.Movie
}
