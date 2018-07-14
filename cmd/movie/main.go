package main

import (
	"flag"
	"time"

	"github.com/jasonlvhit/gocron"

	"github.com/kunhou/TMDB/config"
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/router"

	"github.com/kunhou/TMDB/db"
	movieRepo "github.com/kunhou/TMDB/movie/repository"
	providerRepo "github.com/kunhou/TMDB/provider/repository"
	providerUcase "github.com/kunhou/TMDB/provider/usecase"
)

func main() {
	_localZone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	// migration
	log.Info("Migrate start")
	var drop = flag.Bool("drop", false, "drop all tables")
	var rollback = flag.Int("rollback", 0, "rollback how many steps")
	flag.Parse()
	if *drop {
		db.Drop()
	}
	db.Migrate(*rollback)

	cfg := config.GetConfig()
	mr := movieRepo.NewPGsqlArticleRepository(db.DB)
	pr := providerRepo.NewTMDBRepository(cfg.TMDBToken)
	pu := providerUcase.NewTmdbUsecase(pr, mr)

	ch := pu.CreateBatchStoreTask()
	log.Info("Service Start")
	gocron.ChangeLoc(_localZone)
	go func() {
		s := gocron.NewScheduler()
		s.Every(1).Day().At("04:00").Do(func() {
			log.Info("Start crawler")
			go pu.StartCrawler(ch)
		})
		<-s.Start()
	}()
	router.Setting(pu).Run()
	log.Info("Service Stop")
}
