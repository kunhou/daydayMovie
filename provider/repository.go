package provider

import (
	"fmt"

	"github.com/kunhou/TMDB/models"
)

type ProviderRepository interface {
	GetMovieWithPage(page int) ([]*models.Movie, error)
	GetMovieTotalPages() (int, error)
	GetMovieLastID() (int, error)
	GetMovieDetail(id int) (*models.Movie, error)
	GetPersonLastID() (int, error)
	GetPersonDetail(id int) (*models.Person, error)
}
type APINotFoundError struct {
	Path string
}

func (e APINotFoundError) Error() string {
	return fmt.Sprintf("Not Found %s", e.Path)
}
