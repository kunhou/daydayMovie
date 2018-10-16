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
	ProviderID       uint          `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_movie"`
	Provider         string        `json:"provider" gorm:"type:varchar(127);not null;unique_index:idx_provider_movie"`
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
	People           []Person      `gorm:"many2many:movie_people;association_foreignkey:id;foreignkey:id"`
	CreatedAt        time.Time     `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt        time.Time     `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

type Person struct {
	ID                 uint           `json:"-" gorm:"primary_key"`
	ProviderID         uint           `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_person"`
	Provider           string         `json:"provider" gorm:"type:varchar(127);not null;unique_index:idx_provider_person"`
	Birthday           time.Time      `json:"birthday" gorm:"type:timestamp without time zone;"`
	Name               string         `json:"name" gorm:"type:varchar(255);not null;index"`
	Deathday           time.Time      `json:"deathday" gorm:"type:timestamp without time zone;"`
	Gender             uint8          `json:"gender" gorm:"not null"`
	Biography          string         `json:"biography" gorm:"type:text;not null"`
	Popularity         float32        `json:"popularity"`
	PlaceOfBirth       string         `json:"placeOfBirth" gorm:"type:varchar(255);not null"`
	Adult              bool           `json:"adult" gorm:"not null"`
	ImdbID             string         `json:"imdbID" gorm:"type:varchar(127);not null"`
	Homepage           string         `json:"homepage" gorm:"type:varchar(255);not null"`
	AlsoKnownAs        pq.StringArray `json:"alsoKnownAs,omitempty" gorm:"type:varchar(127)[];not null"`
	ProfilePath        string         `json:"profilePath" gorm:"type:varchar(255);not null"`
	Movies             []Movie        `gorm:"many2many:movie_people;association_foreignkey:id;foreignkey:id"`
	KnownForDepartment string         `json:"knownForDepartment" gorm:"type:varchar(255);not null;default:''"`
	CreatedAt          time.Time      `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt          time.Time      `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

// Set User's table name to be `people`
func (Person) TableName() string {
	return "people"
}

type TV struct {
	ID              int           `json:"id,omitempty" gorm:"primary_key"`
	ProviderID      uint          `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_TV"`
	Provider        string        `json:"-" gorm:"type:varchar(127);not null;unique_index:idx_provider_TV"`
	BackdropPath    string        `json:"backdrop_path" gorm:"type:varchar(255);not null"`
	CreatedByString string        `json:"-" gorm:"type:jsonb;not null"`
	EpisodeRunTime  pq.Int64Array `json:"episode_run_time" gorm:"type:integer[];not null"`
	FirstAirDate    time.Time     `json:"first_air_date"`
	GenreIds        pq.Int64Array `json:"genreIDs" gorm:"type:integer[];not null"`
	Homepage        string        `json:"homepage" gorm:"type:varchar(255);not null"`
	InProduction    bool          `json:"in_production"`
	// Languages      []string      `json:"languages"`
	LastAirDate            time.Time      `json:"last_air_date"`
	LastEpisodeToAirString string         `json:"-" gorm:"type:jsonb;not null"`
	Name                   string         `json:"name"`
	NextEpisodeToAirString string         `json:"-" gorm:"type:jsonb;not null"`
	NetworksString         string         `json:"-" gorm:"type:jsonb;not null"`
	NumberOfEpisodes       int            `json:"number_of_episodes"`
	NumberOfSeasons        int            `json:"number_of_seasons"`
	OriginCountry          pq.StringArray `json:"origin_country" gorm:"type:varchar(127)[];not null"`
	OriginalLanguage       string         `json:"original_language"`
	OriginalName           string         `json:"original_name"`
	Overview               string         `json:"overview"`
	Popularity             float64        `json:"popularity"`
	PosterPath             string         `json:"poster_path"`
	// ProductionCompanies interface{} `json:"production_companies"`
	Seasons     Season    `json:"seasons"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	VoteAverage float64   `json:"vote_average" gorm:"not null;default:'0'"`
	VoteCount   int       `json:"vote_count"  gorm:"not null;default:'0'"`
	CreatedAt   time.Time `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

func (TV) TableName() string {
	return "tv"
}

type Season struct {
	ID           int       `json:"id,omitempty" gorm:"primary_key"`
	TVID         int       `json:"-" gorm:"column:tv_id;not null"`
	AirDate      time.Time `json:"airDate"`
	EpisodeCount int       `json:"episode_count" gorm:"not null"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	Overview     string    `json:"overview" gorm:"type:text;not null"`
	PosterPath   string    `json:"poster_path" gorm:"type:varchar(255);not null"`
	SeasonNumber int       `json:"season_number;not null"`
	VoteAverage  float64   `json:"vote_average" gorm:"not null;default:'0'"`
	VoteCount    int       `json:"vote_count"  gorm:"not null;default:'0'"`
	CreatedAt    time.Time `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

type Credit struct {
	ID         uint   `gorm:"primary_key"`
	PersonID   uint   `gorm:"not null;unique_index:idx_person_cast_type"`
	CastID     uint   `gorm:"not null;unique_index:idx_person_cast_type"`
	Cast       string `gorm:"type:varchar(255);not null;unique_index:idx_person_cast_type"`
	Type       string `gorm:"type:varchar(255);not null;unique_index:idx_person_cast_type"`
	Order      int
	Character  string `gorm:"type:varchar(255)"`
	Department string `gorm:"type:varchar(255)"`
}

func Migrate(rollback int) {
	DB.LogMode(true)
	m := gormigrate.New(DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "201807221000",
			Migrate: func(tx *gorm.DB) error {
				tx.Exec("DROP TABLE if exists people cascade;")
				return tx.AutoMigrate(&Person{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "201807221100",
			Migrate: func(tx *gorm.DB) error {
				return tx.Table("movie_people").
					AddForeignKey("person_id", "people(id)", "CASCADE", "NO ACTION").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "201809151000",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&TV{}, &Season{}).Error; err != nil {
					return err
				}
				if err := tx.Exec(`CREATE UNIQUE INDEX idx_tv_season ON "seasons"(tv_id, season_number)`).Error; err != nil {
					return err
				}
				if err := tx.Model(&Season{}).AddForeignKey("tv_id", "tv(id)", "CASCADE", "NO ACTION").Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "201809291100",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&Person{}).Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "201810162330",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&Credit{}).Error; err != nil {
					return err
				}
				if err := tx.Table("credits").AddForeignKey("person_id", "people(id)", "NO ACTION", "NO ACTION").Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})
	m.InitSchema(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&Movie{}, &Person{}, &TV{}, &Season{}, &Credit{}).Error; err != nil {
			return err
		}
		if err := tx.Table("movie_people").AddForeignKey("movie_id", "movies(id)", "CASCADE", "NO ACTION").
			AddForeignKey("person_id", "people(id)", "CASCADE", "NO ACTION").Error; err != nil {
			return err
		}
		if err := tx.Model(&Season{}).AddForeignKey("tv_id", "tvs(id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}
		if err := tx.Table("credits").AddForeignKey("person_id", "people(id)", "NO ACTION", "NO ACTION").Error; err != nil {
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
