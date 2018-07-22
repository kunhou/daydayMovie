package models

import (
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	ID               uint          `json:"-" gorm:"primary_key"`
	ProviderID       uint          `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_movie"`
	Provider         string        `json:"-" gorm:"type:varchar(127);not null;unique_index:idx_provider_movie"`
	Title            string        `json:"title" gorm:"type:varchar(255);not null;index"`
	OriginalTitle    string        `json:"originalTitle" gorm:"type:varchar(255);not null;index"`
	Popularity       float32       `json:"popularity"`
	VoteAverage      float32       `json:"voteAverage"`
	VoteCount        int           `json:"voteCount"`
	PosterPath       string        `json:"posterPath" gorm:"type:varchar(255);not null"`
	OriginalLanguage string        `json:"-" gorm:"type:varchar(127);not null"`
	GenreIds         pq.Int64Array `json:"genreIDs" gorm:"type:integer[];not null"`
	BackdropPath     string        `json:"backdropPath" gorm:"type:varchar(255);not null"`
	Adult            bool          `json:"adult" gorm:"not null"`
	Overview         string        `json:"overview" gorm:"type:text;not null"`
	ReleaseDate      *time.Time    `json:"releaseDate" gorm:"type:timestamp without time zone"`
	People           []Person      `json:"-" gorm:"many2many:movie_people;association_foreignkey:id;foreignkey:id"`
	CreatedAt        time.Time     `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt        time.Time     `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

type MovieIntro struct {
	ID            uint    `json:"id" gorm:"primary_key"`
	Title         string  `json:"title" gorm:"type:varchar(255);not null;index"`
	OriginalTitle string  `json:"originalTitle" gorm:"type:varchar(255);not null;index"`
	PosterPath    string  `json:"posterPath" gorm:"type:varchar(255);not null"`
	BackdropPath  string  `json:"backdropPath" gorm:"type:varchar(255);not null"`
	Popularity    float32 `json:"popularity"`
}

func (MovieIntro) TableName() string {
	return "movies"
}
