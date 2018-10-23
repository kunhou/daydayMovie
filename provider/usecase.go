package provider

import "github.com/kunhou/TMDB/models"

type ProviderUsecase interface {
	StartCrawlerMovie(ch chan *models.Movie)
	CreateBatchStoreMovieTask() chan *models.Movie
	StartCrawlerPerson(ch chan *models.Person)
	CreateBatchStorePersonTask() chan *models.Person
	StartCrawlerTV(ch chan *models.TV)
	CreateStoreTVTask() chan *models.TV
	StartCrawlerCredit(ch chan *models.Credit, pch chan *models.Person)
	CreateStoreCreditTask() chan *models.Credit
	StartCrawlerPopularMovie(ch chan *models.Credit, pch chan *models.Person, mch chan *models.Movie)
	StartCrawlerPopularTV(ch chan *models.Credit, pch chan *models.Person, tvch chan *models.TV)
}
