package usecase

import (
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/movie"
	"github.com/pkg/errors"
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

func (m *MovieUsecase) MovieList(page, limit int, order map[string]string, query map[string]interface{}) ([]*models.MovieIntro, *models.Page, error) {
	movieIntros, pageInfo, err := m.movieRepos.MovieList(page, limit, order, query)
	if err != nil {
		return nil, nil, err
	}
	movieIDs := []uint{}
	movieMap := map[uint]*models.MovieIntro{}
	for i, m := range movieIntros {
		movieMap[m.ProviderID] = movieIntros[i]
		movieIDs = append(movieIDs, m.ProviderID)
		movieIntros[i].Directing = []models.PersonIntro{}
		movieIntros[i].Cast = []models.PersonIntro{}
	}
	jobType := models.JobDirecting
	directCredits, err := m.movieRepos.CreditIndex(models.CastMovie, &movieIDs, nil, &jobType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get directing credit index")
	}
	dirPeopleIDs := []uint{}
	dirMap := map[uint]*models.MovieIntro{}
	for _, c := range directCredits {
		dirMap[c.PersonID] = movieMap[c.CastID]
		dirPeopleIDs = append(dirPeopleIDs, c.PersonID)
	}
	dirPeople, err := m.movieRepos.PeopleInfoByIDs(dirPeopleIDs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get people intro")
	}
	for i, p := range dirPeople {
		m, ok := dirMap[p.ID]
		if !ok {
			log.Error("not found")
			continue
		}
		m.Directing = append(m.Directing, *dirPeople[i])
	}

	castType := models.CreditTypeCast
	castCredits, err := m.movieRepos.CreditIndex(models.CastMovie, &movieIDs, nil, &castType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get cast credit index")
	}
	castPeopleIDs := []uint{}
	type castMovie struct {
		intro *models.MovieIntro
		order *uint32
	}
	castMap := map[uint]castMovie{}
	for _, c := range castCredits {
		castMap[c.PersonID] = castMovie{
			intro: movieMap[c.CastID],
			order: c.Order,
		}
		castPeopleIDs = append(castPeopleIDs, c.PersonID)
	}
	castPeople, err := m.movieRepos.PeopleInfoByIDs(castPeopleIDs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get cast people intro")
	}
	for i, p := range castPeople {
		m, ok := castMap[p.ID]
		if !ok {
			log.Error("not found")
			continue
		}
		castPeople[i].Order = m.order
		m.intro.Cast = append(m.intro.Cast, *castPeople[i])
	}
	return movieIntros, pageInfo, nil
}

func (m *MovieUsecase) MovieDetail(id uint) (*models.Movie, error) {
	movie, err := m.movieRepos.MovieDetail(id)
	if err != nil {
		return nil, err
	}
	jobType := models.JobDirecting
	movieIDs := []uint{movie.ProviderID}
	dirPeople, err := m.movieRepos.CreditPeople(models.CastMovie, &movieIDs, nil, &jobType)
	if err != nil {
		return nil, errors.Wrap(err, "get directing people index")
	}
	movie.Directing = dirPeople

	castType := models.CreditTypeCast
	castPeople, err := m.movieRepos.CreditPeople(models.CastMovie, &movieIDs, nil, &castType)
	if err != nil {
		return nil, errors.Wrap(err, "get cast people index")
	}
	movie.Cast = castPeople
	return movie, nil
}

func (m *MovieUsecase) TVStore(t *models.TV) (uint, error) {
	return m.movieRepos.TVStore(t)
}

func (m *MovieUsecase) TVList(page, limit int, order map[string]string) ([]*models.TVIntro, *models.Page, error) {
	return m.movieRepos.TVList(page, limit, order)
}

func (m *MovieUsecase) PeopleList(page, limit int, order map[string]string, search map[string]interface{}) ([]*models.Person, *models.Page, error) {
	return m.movieRepos.PeopleList(page, limit, order, search)
}
