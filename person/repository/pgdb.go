package repository

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/person"
)

type pgsqlRepository struct {
	Conn *gorm.DB
}

func NewPGsqlPersonRepository(Conn *gorm.DB) person.PersonRepository {
	return &pgsqlRepository{Conn}
}

func (p *pgsqlRepository) Store(pn *models.Person) (uint, error) {
	if err := p.Conn.Where(models.Person{ProviderID: pn.ProviderID, Provider: pn.Provider}).
		Assign(models.Person{
			Name:         pn.Name,
			Birthday:     pn.Birthday,
			Deathday:     pn.Deathday,
			Gender:       pn.Gender,
			Biography:    pn.Biography,
			Popularity:   pn.Popularity,
			PlaceOfBirth: pn.PlaceOfBirth,
			Adult:        pn.Adult,
			ImdbID:       pn.ImdbID,
			Homepage:     pn.Homepage,
			AlsoKnownAs:  pn.AlsoKnownAs,
			ProfilePath:  pn.ProfilePath,
		}).FirstOrCreate(&pn).Error; err != nil {
		return 0, err
	}
	return pn.ID, nil
}

var TIME_FORMAT = "2006-01-02 15:04:05"

func (p *pgsqlRepository) BatchStore(people []*models.Person) error {
	if len(people) == 0 {
		return nil
	}
	var rows []string
	for _, p := range people {
		birthday := "NULL"
		if p.Birthday != nil {
			birthday = p.Birthday.Format(TIME_FORMAT)
			birthday = fmt.Sprintf("'%s'", birthday)
		}
		deathday := "NULL"
		if p.Deathday != nil {
			deathday = p.Deathday.Format(TIME_FORMAT)
			deathday = fmt.Sprintf("'%s'", deathday)
		}
		alsoKnows, err := p.AlsoKnownAs.Value()
		if err != nil {
			log.WithError(err).Error("GenreIds Parse fail")
		}
		row := fmt.Sprintf("(%d,'%s','%s',%s,%s,'%d','%s','%f','%s','%t','%s','%s','%s', '%s', now(), now())",
			p.ProviderID, p.Provider, strings.Replace(p.Name, "'", "''", -1), birthday, deathday, p.Gender, strings.Replace(p.Biography, "'", "''", -1), p.Popularity, strings.Replace(p.PlaceOfBirth, "'", "''", -1), p.Adult, p.ImdbID, p.Homepage, strings.Replace(alsoKnows.(string), "'", "''", -1), p.ProfilePath)
		rows = append(rows, row)
	}
	sqlStmt := "INSERT INTO people (provider_id,provider,name,birthday,deathday,gender,biography,popularity,place_of_birth,adult,imdb_id,homepage,also_known_as,profile_path,created_at,updated_at) VALUES %s ON CONFLICT (provider, provider_id) DO UPDATE SET " +
		"name = excluded.name, birthday = excluded.birthday, deathday = excluded.deathday, gender = excluded.gender, biography = excluded.biography, popularity = excluded.popularity, place_of_birth = excluded.place_of_birth, adult = excluded.adult, imdb_id = excluded.imdb_id, homepage = excluded.homepage, also_known_as = excluded.also_known_as, profile_path = excluded.profile_path, updated_at = excluded.updated_at;"
	sqlStmt = fmt.Sprintf(sqlStmt, strings.Join(rows, ","))
	if err := p.Conn.Exec(sqlStmt).Error; err != nil {
		return err
	}

	return nil
}
