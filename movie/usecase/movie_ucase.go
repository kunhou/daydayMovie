package usecase

import (
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/movie"
)

type MovieUsecase struct {
	movieRepos movie.MovieRepository
}

func NewMovieUsecase(m movie.MovieRepository) movie.MovieUsecase {
	return &MovieUsecase{
		movieRepos: m,
	}
}

func (m *MovieUsecase) Store(movie *models.Movie) (uint, error) {
	return m.movieRepos.Store(movie)

}
func (m *MovieUsecase) BatchStore(movies []*models.Movie) error {
	return m.movieRepos.BatchStore(movies)
}
