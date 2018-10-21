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

func (m *MovieUsecase) MovieList(page, limit int, order map[string]string) ([]*models.MovieIntro, *models.Page, error) {
	movieIntros, pageInfo, err := m.movieRepos.MovieList(page, limit, order)
	if err != nil {
		return nil, nil, err
	}
	movieIDs := []uint{}
	movieMap := map[uint]*models.MovieIntro{}
	for i, m := range movieIntros {
		movieMap[m.ProviderID] = movieIntros[i]
		movieIDs = append(movieIDs, m.ProviderID)
		movieIntros[i].Directing = []models.PersonIntro{}
	}
	depType := models.DepartmentDirecting
	credits, err := m.movieRepos.CreditIndex(models.CastMovie, &movieIDs, nil, &depType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get credit index")
	}
	peopleIDs := []uint{}
	pMap := map[uint]*models.MovieIntro{}
	for _, c := range credits {
		pMap[c.PersonID] = movieMap[c.CastID]
		peopleIDs = append(peopleIDs, c.PersonID)
	}
	people, err := m.movieRepos.PeopleInfoByIDs(peopleIDs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get people intro")
	}
	for i, p := range people {
		m, ok := pMap[p.ID]
		if !ok {
			log.Error("not found")
			continue
		}
		m.Directing = append(m.Directing, *people[i])
	}
	return movieIntros, pageInfo, nil
}

func (m *MovieUsecase) MovieDetail(id uint) (*models.Movie, error) {
	return m.movieRepos.MovieDetail(id)
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
