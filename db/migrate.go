package db

import (
	"time"

	"github.com/go-gormigrate/gormigrate"
	"github.com/jinzhu/gorm"
	"github.com/kunhou/TMDB/log"
	"github.com/lib/pq"
)

type Movie struct {
	ID               uint          `json:"-" gorm:"primary_key"`
	TmdbID           uint          `json:"-" gorm:"column:tmdb_id;not null;unique_index"`
	Title            string        `json:"title" gorm:"type:varchar(255);not null;index"`
	OriginalTitle    string        `json:"originalTitle" gorm:"type:varchar(255);not null;index"`
	Popularity       float32       `json:"popularity"`
	VoteAverage      float32       `json:"voteAverage"`
	VoteCount        int           `json:"voteCount"`
	PosterPath       string        `json:"posterPath" gorm:"type:varchar(255);not null"`
	OriginalLanguage string        `json:"-" gorm:"type:varchar(127);not null"`
	GenreIds         pq.Int64Array `json:"-" gorm:"type:integer[];not null"`
	BackdropPath     string        `json:"backdropPath" gorm:"type:varchar(255);not null"`
	Adult            bool          `json:"adult" gorm:"not null"`
	Overview         string        `json:"overview" gorm:"type:text;not null"`
	ReleaseDate      time.Time     `json:"releaseDate" gorm:"type:timestamp without time zone"`
	People           []Person      `gorm:"many2many:movie_people;association_foreignkey:tmdbID;foreignkey:tmdbID"`
	CreatedAt        time.Time     `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt        time.Time     `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

type Person struct {
	ID           uint           `json:"-" gorm:"primary_key"`
	TmdbID       uint           `json:"-" gorm:"column:tmdb_id;not null;unique_index"`
	Birthday     time.Time      `json:"birthday" gorm:"type:timestamp without time zone;"`
	Deathday     time.Time      `json:"deathday" gorm:"type:timestamp without time zone;"`
	Gender       uint8          `json:"gender" gorm:"not null"`
	Biography    string         `json:"biography" gorm:"type:text;not null"`
	Popularity   float32        `json:"popularity"`
	PlaceOfBirth string         `json:"placeOfBirth" gorm:"type:varchar(255);not null"`
	Adult        bool           `json:"adult" gorm:"not null"`
	ImdbID       string         `json:"imdbID" gorm:"type:varchar(127);not null"`
	Homepage     string         `json:"homepage" gorm:"type:varchar(255);not null"`
	AlsoKnownAs  pq.StringArray `json:"alsoKnownAs,omitempty" gorm:"type:varchar(127)[];not null"` // {"HTTP", "HTTPS"}
	Movies       []Movie        `gorm:"many2many:movie_people;association_foreignkey:tmdbID;foreignkey:tmdbID"`
	CreatedAt    time.Time      `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt    time.Time      `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

// Set User's table name to be `people`
func (Person) TableName() string {
	return "people"
}

func Migrate(rollback int) {
	DB.LogMode(true)
	m := gormigrate.New(DB, gormigrate.DefaultOptions, []*gormigrate.Migration{})
	m.InitSchema(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&Movie{}, &Person{}).Error; err != nil {
			return err
		}
		if err := tx.Table("movie_people").AddForeignKey("movie_tmdb_id", "movies(tmdb_id)", "CASCADE", "NO ACTION").
			AddForeignKey("person_tmdb_id", "people(tmdb_id)", "CASCADE", "NO ACTION").Error; err != nil {
			return err
		}
		return nil
	})
	if err := m.Migrate(); err != nil {
		log.WithError(err).Error("Database Migration Failed")
	}
	log.Info("Migrate Finished")
	DB.LogMode(false)
}

func Drop() {
	DB.LogMode(true)
	log.Info("Drop Tables...")
	// drop tables
	DB.DropTableIfExists(&gormigrate.Migration{})
	DB.DropTableIfExists(&Movie{})
	DB.DropTableIfExists(&Person{})
	DB.LogMode(false)
}
