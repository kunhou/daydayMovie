package main

import (
	"flag"

	"github.com/kunhou/TMDB/config"
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/router"

	"github.com/kunhou/TMDB/db"
	movieRepo "github.com/kunhou/TMDB/movie/repository"
	movieUcase "github.com/kunhou/TMDB/movie/usecase"
	personRepo "github.com/kunhou/TMDB/person/repository"
	providerRepo "github.com/kunhou/TMDB/provider/repository"
	providerUcase "github.com/kunhou/TMDB/provider/usecase"
)

func main() {
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
	mr := movieRepo.NewPGsqlMovieRepository(db.DB)
	personr := personRepo.NewPGsqlPersonRepository(db.DB)
	pr := providerRepo.NewTMDBRepository(cfg.TMDBToken)
	pu := providerUcase.NewTmdbUsecase(pr, mr, personr)
	mu := movieUcase.NewMovieUsecase(mr)

	log.Info("Service Start")
	router.Setting(pu, mu).Run()
	log.Info("Service Stop")
}
