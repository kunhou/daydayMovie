package provider

import "github.com/kunhou/TMDB/models"

type ProviderRepository interface {
	GetMovieWithPage(page int) ([]*models.Movie, error)
	GetMovieTotalPages() (int, error)
}
