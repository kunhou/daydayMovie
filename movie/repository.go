package movie

import "github.com/kunhou/TMDB/models"

type MovieRepository interface {
	Store(m *models.Movie) (uint, error)
	BatchStore(movies []*models.Movie) error
}
