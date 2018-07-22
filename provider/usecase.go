package provider

import "github.com/kunhou/TMDB/models"

type ProviderUsecase interface {
	StartCrawlerMovie(ch chan *models.Movie)
	CreateBatchStoreMovieTask() chan *models.Movie
	StartCrawlerPerson(ch chan *models.Person)
	CreateBatchStorePersonTask() chan *models.Person
}
