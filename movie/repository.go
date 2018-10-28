package movie

import (
	"github.com/kunhou/TMDB/models"
)

type MovieRepository interface {
	Store(m *models.Movie) (uint, error)
	BatchStore(movies []*models.Movie) error
	MovieList(page, limit int, order map[string]string, query map[string]interface{}) ([]*models.MovieIntro, *models.Page, error)
	MovieDetail(id uint) (*models.Movie, error)
	TVStore(t *models.TV) (uint, error)
	TVList(page, limit int, order map[string]string) ([]*models.TVIntro, *models.Page, error)
	PeopleList(page, limit int, order map[string]string, search map[string]interface{}) ([]*models.Person, *models.Page, error)
	PeopleInfoByIDs(pIDs []uint) ([]*models.PersonIntro, error)
	PeopleIDByProviderID(pIDs uint) (uint, error)
	CreditStore(c *models.Credit) (uint, error)
	CreditIndex(castType string, castIDs *[]uint, peopleIDs *[]uint, department *string) ([]*models.Credit, error)
	CreditPeople(castType string, castIDs *[]uint, peopleIDs *[]uint, job *string) ([]models.PersonIntro, error)
}
