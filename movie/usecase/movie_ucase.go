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

func (m *MovieUsecase) MovieList(page, limit int, order map[string]string) ([]*models.MovieIntro, *models.Page, error) {
	return m.movieRepos.MovieList(page, limit, order)
}

func (m *MovieUsecase) MovieDetail(id uint) (*models.Movie, error) {
	return m.movieRepos.MovieDetail(id)
}

func (m *MovieUsecase) TVStore(t *models.TV) (uint, error) {
	return m.movieRepos.TVStore(t)
}

func (m *MovieUsecase) TVList(page, limit int, order map[string]string) ([]*models.TVIntro, *models.Page, error) {
	return m.movieRepos.TVList(page, limit, order)
}
