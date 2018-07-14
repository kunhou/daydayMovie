package usecase

import (
	"time"

	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/movie"
	"github.com/kunhou/TMDB/provider"
)

type tmdbUsecase struct {
	providerRepo provider.ProviderRepository
	movieRepo    movie.MovieRepository
}

var movieWriteChan chan *models.Movie

func NewTmdbUsecase(repo provider.ProviderRepository, movieRepo movie.MovieRepository) provider.ProviderUsecase {
	return &tmdbUsecase{repo, movieRepo}
}

func (tu *tmdbUsecase) StartCrawler(ch chan *models.Movie) {
	total, err := tu.providerRepo.GetMovieTotalPages()
	if err != nil {
		log.WithError(err).Error("Get discover Fail")
	}
	for page := 1; page <= total; page++ {
		time.Sleep(400 * time.Millisecond)
		go func(p int) {
			movies, err := tu.providerRepo.GetMovieWithPage(p)
			if err != nil {
				log.WithError(err).Error("Get discover Fail")
			}
			for _, m := range movies {
				ch <- m
			}
		}(page)
	}
	return
}

func (tu *tmdbUsecase) CreateBatchStoreTask() chan *models.Movie {
	if movieWriteChan != nil {
		return movieWriteChan
	}
	movieWriteChan = make(chan *models.Movie, 100000)
	movieBuff := []*models.Movie{}
	go func() {
		t := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-t.C:
				if len(movieBuff) == 0 {
					continue
				}
				if err := tu.movieRepo.BatchStore(movieBuff); err != nil {
					log.WithError(err).Error("Movie Task store fail")
				}
				movieBuff = []*models.Movie{}
			case movie := <-movieWriteChan:
				movieBuff = append(movieBuff, movie)
				if len(movieBuff) > 0 {
					if err := tu.movieRepo.BatchStore(movieBuff); err != nil {
						log.WithError(err).Error("Movie Task store fail")
					}
					movieBuff = []*models.Movie{}
				}
			}
		}
	}()

	return movieWriteChan
}
