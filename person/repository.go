package person

import "github.com/kunhou/TMDB/models"

type PersonRepository interface {
	Store(person *models.Person) (uint, error)
	BatchStore(people []*models.Person) error
}
