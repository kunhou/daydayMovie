package movie

import "github.com/kunhou/TMDB/models"

type MovieUsecase interface {
	Store(m *models.Movie) (uint, error)
	BatchStore(movies []*models.Movie) error
}
