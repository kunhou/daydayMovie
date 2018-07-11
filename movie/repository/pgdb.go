package repository

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/movie"
)

type pgsqlRepository struct {
	Conn *gorm.DB
}

func NewPGsqlArticleRepository(Conn *gorm.DB) movie.MovieRepository {
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
