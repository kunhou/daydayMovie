package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type CreatedBy struct {
	ID          int    `json:"id"`
	CreditID    string `json:"credit_id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	ProfilePath string `json:"profile_path"`
}
type LastEpisodeToAir struct {
	AirDate        string  `json:"airDate"`
	EpisodeNumber  int     `json:"episode_number"`
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	ShowID         int     `json:"show_id"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}
type NextEpisodeToAir struct {
	ID             int     `json:"id"`
	AirDate        string  `json:"airDate"`
	EpisodeNumber  int     `json:"episode_number"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	ShowID         int     `json:"show_id"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}
type Networks struct {
	Name          string `json:"name"`
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}
type TV struct {
	ID              uint          `json:"id,omitempty" gorm:"primary_key"`
	ProviderID      uint          `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_movie"`
	Provider        string        `json:"-" gorm:"type:varchar(127);not null;unique_index:idx_provider_movie"`
	BackdropPath    string        `json:"backdrop_path" gorm:"type:varchar(255);not null"`
	CreatedBy       []CreatedBy   `json:"created_by"`
	CreatedByString string        `json:"-" gorm:"type:jsonb;not null"`
	EpisodeRunTime  pq.Int64Array `json:"episode_run_time"`
	FirstAirDate    *time.Time    `json:"first_air_date"`
	GenreIds        pq.Int64Array `json:"genreIDs" gorm:"type:integer[];not null"`
	Homepage        string        `json:"homepage" gorm:"type:varchar(255);not null"`
	InProduction    bool          `json:"in_production"`
	// Languages      []string      `json:"languages"`
	LastAirDate            *time.Time       `json:"last_air_date"`
	LastEpisodeToAir       LastEpisodeToAir `json:"last_episode_to_air"`
	LastEpisodeToAirString string           `json:"-" gorm:"type:jsonb;not null"`
	Name                   string           `json:"name"`
	NextEpisodeToAir       NextEpisodeToAir `json:"next_episode_to_air"`
	NextEpisodeToAirString string           `json:"-" gorm:"type:jsonb;not null"`
	Networks               []Networks       `json:"networks"`
	NetworksString         string           `json:"-" gorm:"type:jsonb;not null"`
	NumberOfEpisodes       int              `json:"number_of_episodes"`
	NumberOfSeasons        int              `json:"number_of_seasons"`
	OriginCountry          pq.StringArray   `json:"origin_country"`
	OriginalLanguage       string           `json:"original_language"`
	OriginalName           string           `json:"original_name"`
	Overview               string           `json:"overview"`
	Popularity             float64          `json:"popularity"`
	PosterPath             string           `json:"poster_path"`
	// ProductionCompanies interface{} `json:"production_companies"`
	Seasons     []Season  `json:"seasons" gorm:"-"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	VoteAverage float64   `json:"vote_average"`
	VoteCount   int       `json:"vote_count"`
	CreatedAt   time.Time `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

func (TV) TableName() string {
	return "tv"
}

func (v *TV) BeforeSave() (err error) {
	createdByStr, err := json.Marshal(v.CreatedBy)
	if err != nil {
		return err
	}
	v.CreatedByString = string(createdByStr)

	lastEpisodeToAirStr, err := json.Marshal(v.LastEpisodeToAir)
	if err != nil {
		return err
	}
	v.LastEpisodeToAirString = string(lastEpisodeToAirStr)

	nextEpisodeToAir, err := json.Marshal(v.NextEpisodeToAir)
	if err != nil {
		return err
	}
	v.NextEpisodeToAirString = string(nextEpisodeToAir)

	networks, err := json.Marshal(v.Networks)
	if err != nil {
		return err
	}
	v.NetworksString = string(networks)

	return
}

func (v *TV) AfterFind(scope *gorm.Scope) (err error) {
	if v.CreatedByString != "" {
		if err := json.Unmarshal([]byte(v.CreatedByString), &v.CreatedBy); err != nil {
			return err
		}
	}
	if v.LastEpisodeToAirString != "" {
		if err := json.Unmarshal([]byte(v.LastEpisodeToAirString), &v.LastEpisodeToAir); err != nil {
			return err
		}
	}
	if v.NextEpisodeToAirString != "" {
		if err := json.Unmarshal([]byte(v.NextEpisodeToAirString), &v.NextEpisodeToAir); err != nil {
			return err
		}
	}
	if v.NetworksString != "" {
		if err := json.Unmarshal([]byte(v.NetworksString), &v.Networks); err != nil {
			return err
		}
	}

	return
}

type Season struct {
	ID           uint       `json:"id,omitempty" gorm:"primary_key"`
	TVID         uint       `json:"-" gorm:"column:tv_id;not null;unique_index:idx_tv_season"`
	AirDate      *time.Time `json:"airDate"`
	EpisodeCount int        `json:"episode_count"`
	Name         string     `json:"name"`
	Overview     string     `json:"overview"`
	PosterPath   string     `json:"poster_path"`
	SeasonNumber int        `json:"season_number;not null;unique_index:idx_tv_season"`
	VoteAverage  float64    `json:"vote_average" gorm:"not null;default:'0'"`
	VoteCount    int        `json:"vote_count"  gorm:"not null;default:'0'"`
	CreatedAt    time.Time  `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt    time.Time  `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

type TVIntro struct {
	ID            uint       `json:"id" gorm:"primary_key"`
	Title         string     `json:"title" gorm:"column:name"`
	Overview      string     `json:"overview"`
	OriginalTitle string     `json:"original_title" gorm:"column:original_name"`
	PosterPath    string     `json:"posterPath" gorm:"type:varchar(255);not null"`
	BackdropPath  string     `json:"backdropPath" gorm:"type:varchar(255);not null"`
	Popularity    float32    `json:"popularity"`
	VoteAverage   float64    `json:"vote_average"`
	VoteCount     int        `json:"vote_count"`
	LastAirDate   *time.Time `json:"last_air_date"`
}

func (TVIntro) TableName() string {
	return "tv"
}
