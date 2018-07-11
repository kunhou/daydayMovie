package provider

import (
	"github.com/kunhou/TMDB/config"
)

var (
	cfg           = config.GetConfig()
	apiURL        = "https://api.themoviedb.org/3/"
	DISCOVER_PATH = "/discover/movie"
)

// func CrawlerRun() {
// 	mr := movieRepo.NewPGsqlArticleRepository(db.DB)
// 	mu := movieUsecase.NewMovieUsecase(mr)
// 	ch := mu.CreateBatchStoreTask()

// 	info, err := getDiscoverMovieByPage(1)
// 	if err != nil {
// 		log.WithError(err).Error("Get discover Fail")
// 	}
// 	for page := 1; page <= info.TotalPages; page++ {
// 		time.Sleep(400 * time.Millisecond)
// 		go func(p int) {
// 			r, err := getDiscoverMovieByPage(p)
// 			if err != nil {
// 				log.WithError(err).Error("Get discover Fail")
// 			}
// 			for _, m := range r.Movies {
// 				ch <- m
// 			}
// 		}(page)
// 	}
// 	return
// }
