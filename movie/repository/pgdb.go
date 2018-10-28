package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/movie"
	"github.com/kunhou/TMDB/utils"
	"github.com/pkg/errors"
)

type pgsqlRepository struct {
	Conn *gorm.DB
}

func NewPGsqlMovieRepository(Conn *gorm.DB) movie.MovieRepository {
	return &pgsqlRepository{Conn}
}

func (p *pgsqlRepository) Store(m *models.Movie) (uint, error) {
	if err := p.Conn.Where(models.Movie{ProviderID: m.ProviderID, Provider: m.Provider}).
		Assign(models.Movie{
			Title:            m.Title,
			OriginalTitle:    m.OriginalTitle,
			Popularity:       m.Popularity,
			VoteAverage:      m.VoteAverage,
			VoteCount:        m.VoteCount,
			PosterPath:       m.PosterPath,
			OriginalLanguage: m.OriginalLanguage,
			GenreIds:         m.GenreIds,
			BackdropPath:     m.BackdropPath,
			Adult:            m.Adult,
			Overview:         m.Overview,
			ReleaseDate:      m.ReleaseDate,
		}).FirstOrCreate(&m).Error; err != nil {
		return 0, err
	}
	return m.ID, nil
}

var TIME_FORMAT = "2006-01-02 15:04:05"

func (p *pgsqlRepository) BatchStore(movies []*models.Movie) error {
	if len(movies) == 0 {
		return nil
	}
	var rows []string
	for _, m := range movies {
		genreIDs, err := m.GenreIds.Value()
		if err != nil {
			log.WithError(err).Error("GenreIds Parse fail")
		}
		releaseDate := "NULL"
		if m.ReleaseDate != nil {
			releaseDate = m.ReleaseDate.Format(TIME_FORMAT)
			releaseDate = fmt.Sprintf("'%s'", releaseDate)
		}
		row := fmt.Sprintf("(%d,'%s','%s','%s','%f','%f','%d','%s','%s','%s','%s','%t','%s', %s, now(), now())",
			m.ProviderID, m.Provider, strings.Replace(m.Title, "'", "''", -1), strings.Replace(m.OriginalTitle, "'", "''", -1), m.Popularity, m.VoteAverage, m.VoteCount, m.PosterPath, m.OriginalLanguage, genreIDs, m.BackdropPath, m.Adult, strings.Replace(m.Overview, "'", "''", -1), releaseDate)
		rows = append(rows, row)
	}
	sqlStmt := "INSERT INTO movies (provider_id, provider, title, original_title, popularity, vote_average, vote_count, poster_path, original_language, genre_ids, backdrop_path, adult, overview, release_date, created_at, updated_at) VALUES %s ON CONFLICT (provider, provider_id) DO UPDATE SET " +
		"title = excluded.title, original_title = excluded.original_title, popularity = excluded.popularity, vote_average = excluded.vote_average, vote_count = excluded.vote_count, poster_path = excluded.poster_path, original_language = excluded.original_language, genre_ids = excluded.genre_ids, backdrop_path = excluded.backdrop_path, adult = excluded.adult, overview = excluded.overview, release_date = excluded.release_date, updated_at = excluded.updated_at;"
	sqlStmt = fmt.Sprintf(sqlStmt, strings.Join(rows, ","))
	if err := p.Conn.Exec(sqlStmt).Error; err != nil {
		return err
	}

	return nil
}

var movieOrder = map[string]bool{
	"id":          true,
	"popularity":  true,
	"releaseDate": true,
}

func (p *pgsqlRepository) MovieList(page, limit int, order map[string]string, query map[string]interface{}) ([]*models.MovieIntro, *models.Page, error) {
	if limit == 0 {
		limit = 20
	}
	if page == 0 {
		page = 1
	}
	db := p.Conn
	oColumn, oType := "popularity", "desc"
	for column, orderType := range order {
		if _, ok := movieOrder[column]; !ok {
			continue
		}
		if !utils.ValidOrderType(orderType) {
			continue
		}
		if strings.EqualFold(column, "releaseDate") {
			oColumn = "release_date"
			db = db.Where("release_date IS NOT NULL")
		} else {
			oColumn = column
		}
		oType = orderType
		break
	}
	if gs, ok := query["genres"]; ok {
		genres := gs.([]string)
		gIDs := []string{}
		for _, gName := range genres {
			for id, name := range models.MovieGenres {
				if strings.EqualFold(gName, name) {
					gIDs = append(gIDs, strconv.Itoa(int(id)))
					break
				}
			}
		}
		db = db.Where("genre_ids @> ARRAY[" + strings.Join(gIDs, ",") + "]")
	}
	db = db.Order(oColumn + " " + oType)

	offset := (page - 1) * limit
	movies := []*models.MovieIntro{}
	var count uint
	if err := db.Table("movies").Count(&count).Error; err != nil {
		return movies, nil, err
	}
	if err := db.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return movies, nil, err
	}

	totalPages := count/uint(limit) + 1
	pages := models.Page{
		TotalPages:   uint(totalPages),
		TotalResults: count,
		Page:         uint(page),
	}
	return movies, &pages, nil
}

func (p *pgsqlRepository) MovieDetail(id uint) (*models.Movie, error) {
	m := models.Movie{
		ID: id,
	}
	if err := p.Conn.Find(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (p *pgsqlRepository) TVStore(t *models.TV) (uint, error) {
	if err := p.Conn.Where(models.TV{ProviderID: t.ProviderID, Provider: t.Provider}).
		Assign(models.TV{
			BackdropPath:     t.BackdropPath,
			CreatedBy:        t.CreatedBy,
			EpisodeRunTime:   t.EpisodeRunTime,
			FirstAirDate:     t.FirstAirDate,
			GenreIds:         t.GenreIds,
			Homepage:         t.Homepage,
			InProduction:     t.InProduction,
			LastAirDate:      t.LastAirDate,
			LastEpisodeToAir: t.LastEpisodeToAir,
			Name:             t.Name,
			NextEpisodeToAir: t.NextEpisodeToAir,
			Networks:         t.Networks,
			NumberOfEpisodes: t.NumberOfEpisodes,
			NumberOfSeasons:  t.NumberOfSeasons,
			OriginCountry:    t.OriginCountry,
			OriginalLanguage: t.OriginalLanguage,
			OriginalName:     t.OriginalName,
			Overview:         t.Overview,
			Popularity:       t.Popularity,
			PosterPath:       t.PosterPath,
			// Seasons:          t.Seasons,
			Status:      t.Status,
			Type:        t.Type,
			VoteAverage: t.VoteAverage,
			VoteCount:   t.VoteCount,
		}).FirstOrCreate(&t).Error; err != nil {
		return 0, errors.Wrap(err, "Store tv Fail")
	}
	for _, s := range t.Seasons {
		if err := p.Conn.Where(models.Season{TVID: t.ID, SeasonNumber: s.SeasonNumber}).
			Assign(models.Season{
				AirDate:      s.AirDate,
				EpisodeCount: s.EpisodeCount,
				Name:         s.Name,
				PosterPath:   s.PosterPath,
				SeasonNumber: s.SeasonNumber,
				VoteAverage:  s.VoteAverage,
				VoteCount:    s.VoteCount,
			}).FirstOrCreate(&s).Error; err != nil {
			return 0, errors.Wrap(err, "Store Season Fail")
		}
	}
	return t.ID, nil
}

func (p *pgsqlRepository) TVList(page, limit int, order map[string]string) ([]*models.TVIntro, *models.Page, error) {
	if limit == 0 {
		limit = 20
	}
	if page == 0 {
		page = 1
	}
	db := p.Conn
	if o, ok := order["popularity"]; ok && utils.ValidOrderType(o) {
		db = db.Order("popularity " + o)
	}
	offset := (page - 1) * limit
	tvs := []*models.TVIntro{}
	var count uint
	if err := db.Table("tv").Count(&count).Error; err != nil {
		return tvs, nil, err
	}
	if err := db.Offset(offset).Limit(limit).Find(&tvs).Error; err != nil {
		return tvs, nil, err
	}

	totalPages := count/uint(limit) + 1
	pages := models.Page{
		TotalPages:   uint(totalPages),
		TotalResults: count,
		Page:         uint(page),
	}
	return tvs, &pages, nil
}

func (p *pgsqlRepository) PeopleList(page, limit int, order map[string]string, search map[string]interface{}) ([]*models.Person, *models.Page, error) {
	if limit == 0 {
		limit = 20
	}
	if page == 0 {
		page = 1
	}
	db := p.Conn
	// if o, ok := order["popularity"]; ok && utils.ValidOrderType(o) {
	// 	db = db.Order("popularity " + o)
	// }
	offset := (page - 1) * limit
	people := []*models.Person{}
	var count uint
	// date_part('month', birthday) = 4 AND date_part('day', birthday) = 4
	if b, ok := search["birthday"]; ok {
		bTime := b.(time.Time)
		db = db.Where("date_part('month', birthday) = ? AND date_part('day', birthday) = ?", bTime.Month(), bTime.Day())
		delete(search, "birthday")
	}
	if err := db.Table("people").Where(search).Count(&count).Error; err != nil {
		return people, nil, err
	}
	if err := db.Offset(offset).Where(search).Limit(limit).Find(&people).Error; err != nil {
		return people, nil, err
	}

	totalPages := count/uint(limit) + 1
	pages := models.Page{
		TotalPages:   uint(totalPages),
		TotalResults: count,
		Page:         uint(page),
	}
	return people, &pages, nil
}

func (p *pgsqlRepository) PeopleInfoByIDs(pIDs []uint) ([]*models.PersonIntro, error) {
	db := p.Conn
	peopleIntros := []*models.PersonIntro{}
	if err := db.Where("id IN (?)", pIDs).Find(&peopleIntros).Error; err != nil && err != gorm.ErrRecordNotFound {
		return peopleIntros, err
	}
	return peopleIntros, nil
}

func (p *pgsqlRepository) PeopleIDByProviderID(pID uint) (uint, error) {
	db := p.Conn
	peopleIntro := models.PersonIntro{}
	if err := db.Where("provider_id = ?", pID).First(&peopleIntro).Error; err != nil {
		return 0, err
	}
	return peopleIntro.ID, nil
}

func (p *pgsqlRepository) CreditIndex(castType string, castIDs *[]uint, peopleIDs *[]uint, job *string) ([]*models.Credit, error) {
	db := p.Conn
	credits := []*models.Credit{}
	db = db.Where("\"cast\" = ?", castType)
	if castIDs != nil {
		ids := []int{}
		for _, i := range *castIDs {
			ids = append(ids, int(i))
		}
		db = db.Where("cast_id IN (?)", ids)
	}
	if peopleIDs != nil {
		ids := []int{}
		for _, i := range *peopleIDs {
			ids = append(ids, int(i))
		}
		db = db.Where("person_id IN (?)", ids)
	}
	if job != nil {
		if strings.EqualFold(*job, models.CreditTypeCast) {
			db = db.Where("type = ?", job).Order("\"order\" ASC")
		} else {
			db = db.Where("job = ?", job)
		}
	}
	if err := db.Find(&credits).Error; err != nil && err != gorm.ErrRecordNotFound {
		return credits, err
	}
	return credits, nil
}

func (p *pgsqlRepository) CreditStore(c *models.Credit) (uint, error) {
	if err := p.Conn.Where(models.Credit{Cast: c.Cast, CastID: c.CastID, PersonID: c.PersonID, Type: c.Type}).
		Assign(models.Credit{
			Order:      c.Order,
			Department: c.Department,
			Job:        c.Job,
		}).FirstOrCreate(&c).Error; err != nil {
		return 0, errors.Wrap(err, "Store credit Fail")
	}
	return c.ID, nil
}

func (p *pgsqlRepository) CreditPeople(castType string, castIDs *[]uint, peopleIDs *[]uint, job *string) ([]models.PersonIntro, error) {
	people := []models.PersonIntro{}
	db := p.Conn.Table("credits").
		Select([]string{"people.*", "credits.order"}).
		Joins("LEFT JOIN people ON person_id = people.id")
	if castIDs != nil {
		ids := []int{}
		for _, i := range *castIDs {
			ids = append(ids, int(i))
		}
		db = db.Where("cast_id IN (?)", ids)
	}
	if peopleIDs != nil {
		ids := []int{}
		for _, i := range *peopleIDs {
			ids = append(ids, int(i))
		}
		db = db.Where("person_id IN (?)", ids)
	}
	if job != nil {
		if strings.EqualFold(*job, models.CreditTypeCast) {
			db = db.Where("type = ?", job).Order("\"order\" ASC")
		} else {
			db = db.Where("job = ?", job)
		}
	}
	if err := db.Find(&people).Error; err != nil && err != gorm.ErrRecordNotFound {
		return people, err
	}
	return people, nil
}
