package person

import "github.com/kunhou/TMDB/models"

type PersonUsecase interface {
	Store(p *models.Person) (uint, error)
	BatchStore(persons []*models.Person) error
}
