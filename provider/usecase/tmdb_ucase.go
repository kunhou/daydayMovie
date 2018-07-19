package usecase

import (
	"sync"
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
var movieWriteSyncOnce sync.Once

func NewTmdbUsecase(repo provider.ProviderRepository, movieRepo movie.MovieRepository) provider.ProviderUsecase {
	return &tmdbUsecase{repo, movieRepo}
}

func (tu *tmdbUsecase) StartCrawler(ch chan *models.Movie) {
	lastestID, err := tu.providerRepo.GetMovieLastID()
	if err != nil {
		log.WithError(err).Error("Get LastID Fail")
	}
	for id := 1; id <= lastestID; id++ {
		time.Sleep(700 * time.Millisecond)
		go func(p int) {
			m, err := tu.providerRepo.GetMovieDetail(p)
			if err != nil {
				if _, ok := err.(provider.APINotFoundError); !ok {
					log.WithError(err).Error("Get discover Fail")
					return
				}
				return
			}
			ch <- m
		}(id)
	}
	return
}

func (tu *tmdbUsecase) CreateBatchStoreTask() chan *models.Movie {
	movieWriteSyncOnce.Do(func() {
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
					if len(movieBuff) > 100 {
						if err := tu.movieRepo.BatchStore(movieBuff); err != nil {
							log.WithError(err).Error("Movie Task store fail")
						}
						movieBuff = []*models.Movie{}
					}
				}
			}
		}()
	})

	return movieWriteChan
}
