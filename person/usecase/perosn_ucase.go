package usecase

import (
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/person"
)

type PersonUsecase struct {
	PersonRepos person.PersonRepository
}

func NewPersonUsecase(p Person.PersonRepository) Person.PersonUsecase {
	return &PersonUsecase{
		PersonRepos: p,
	}
}

func (p *PersonUsecase) Store(person *models.Person) (uint, error) {
	return p.PersonRepos.Store(person)

}
func (p *PersonUsecase) BatchStore(people []*models.Person) error {
	return p.PersonRepos.BatchStore(people)
}
