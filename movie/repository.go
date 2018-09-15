package movie

import "github.com/kunhou/TMDB/models"

type MovieRepository interface {
	Store(m *models.Movie) (uint, error)
	BatchStore(movies []*models.Movie) error
	MovieList(page, limit int, order map[string]string) ([]*models.MovieIntro, *models.Page, error)
	MovieDetail(id uint) (*models.Movie, error)
	TVStore(t *models.TV) (uint, error)
}
