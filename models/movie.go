package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

const TIME_FORMAT = "2006-01-02"

var MovieGenres = map[int64]string{
	28:    "動作",
	12:    "冒險",
	16:    "動畫",
	35:    "喜劇",
	80:    "犯罪",
	99:    "紀錄",
	18:    "劇情",
	10751: "家庭",
	14:    "奇幻",
	36:    "歷史",
	27:    "恐怖",
	10402: "音樂",
	9648:  "懸疑",
	10749: "愛情",
	878:   "科幻",
	10770: "電視電影",
	53:    "驚悚",
	10752: "戰爭",
	37:    "西部",
}

type Movie struct {
	ID               uint          `json:"id" gorm:"primary_key"`
	ProviderID       uint          `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_movie"`
	Provider         string        `json:"-" gorm:"type:varchar(127);not null;unique_index:idx_provider_movie"`
	Title            string        `json:"title" gorm:"type:varchar(255);not null;index"`
	OriginalTitle    string        `json:"originalTitle" gorm:"type:varchar(255);not null;index"`
	Popularity       float32       `json:"popularity"`
	VoteAverage      float32       `json:"voteAverage"`
	VoteCount        int           `json:"voteCount"`
	PosterPath       string        `json:"posterPath" gorm:"type:varchar(255);not null"`
	OriginalLanguage string        `json:"-" gorm:"type:varchar(127);not null"`
	GenreIds         pq.Int64Array `json:"-" gorm:"type:integer[];not null"`
	Genres           []string      `json:"genres" gorm:"-"`
	BackdropPath     string        `json:"backdropPath" gorm:"type:varchar(255);not null"`
	Adult            bool          `json:"adult" gorm:"not null"`
	Overview         string        `json:"overview" gorm:"type:text;not null"`
	ReleaseDate      *time.Time    `json:"-" gorm:"type:timestamp without time zone"`
	Directing        []PersonIntro `json:"directing" gorm:"-"`
	Cast             []PersonIntro `json:"cast" gorm:"-"`
	CreatedAt        time.Time     `json:"-" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt        time.Time     `json:"-" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

func (m *Movie) AfterFind(scope *gorm.Scope) (err error) {
	m.Genres = []string{}
	for _, id := range m.GenreIds {
		m.Genres = append(m.Genres, MovieGenres[id])
	}
	return
}

func (m *Movie) MarshalJSON() ([]byte, error) {
	type Alias Movie
	var releaseDate *string
	if m.ReleaseDate != nil {
		r := m.ReleaseDate.Format(TIME_FORMAT)
		releaseDate = &r
	}
	return json.Marshal(&struct {
		*Alias
		ReleaseDate *string `json:"releaseDate"`
	}{
		Alias:       (*Alias)(m),
		ReleaseDate: releaseDate,
	})
}

type MovieIntro struct {
	ID            uint          `json:"id" gorm:"primary_key"`
	ProviderID    uint          `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_movie"`
	Title         string        `json:"title" gorm:"type:varchar(255);not null;index"`
	OriginalTitle string        `json:"originalTitle" gorm:"type:varchar(255);not null;index"`
	PosterPath    string        `json:"posterPath" gorm:"type:varchar(255);not null"`
	BackdropPath  string        `json:"backdropPath" gorm:"type:varchar(255);not null"`
	Popularity    float32       `json:"popularity"`
	VoteAverage   float64       `json:"vote_average"`
	VoteCount     int           `json:"vote_count"`
	Overview      string        `json:"overview" gorm:"type:text;not null"`
	ReleaseDate   *time.Time    `json:"releaseDate" gorm:"type:timestamp without time zone"`
	Directing     []PersonIntro `json:"directing" gorm:"-"`
	Cast          []PersonIntro `json:"cast" gorm:"-"`
	GenreIds      pq.Int64Array `json:"-" gorm:"type:integer[];not null"`
	Genres        []string      `json:"genres" gorm:"-"`
}

func (MovieIntro) TableName() string {
	return "movies"
}

func (m *MovieIntro) AfterFind(scope *gorm.Scope) (err error) {
	m.Genres = []string{}
	for _, id := range m.GenreIds {
		m.Genres = append(m.Genres, MovieGenres[id])
	}
	return
}

func (m *MovieIntro) MarshalJSON() ([]byte, error) {
	type Alias MovieIntro
	var releaseDate *string
	if m.ReleaseDate != nil {
		r := m.ReleaseDate.Format(TIME_FORMAT)
		releaseDate = &r
	}
	return json.Marshal(&struct {
		*Alias
		ReleaseDate *string `json:"releaseDate"`
	}{
		Alias:       (*Alias)(m),
		ReleaseDate: releaseDate,
	})
}
