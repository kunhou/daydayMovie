package models

import (
	"time"

	"github.com/lib/pq"
)

type Person struct {
	ID                 uint           `json:"-" gorm:"primary_key"`
	ProviderID         uint           `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_person"`
	Provider           string         `json:"provider" gorm:"type:varchar(127);not null;unique_index:idx_provider_person"`
	Name               string         `json:"name" gorm:"type:varchar(255);not null;index"`
	Birthday           *time.Time     `json:"birthday" gorm:"type:timestamp without time zone;"`
	Deathday           *time.Time     `json:"deathday" gorm:"type:timestamp without time zone;"`
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
	Order              *uint32        `json:"order,omitempty"`
	CreatedAt          time.Time      `json:"createdAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
	UpdatedAt          time.Time      `json:"updatedAt,omitempty" gorm:"type:timestamp without time zone;not null;default:'now()'"`
}

func (Person) TableName() string {
	return "people"
}

type PersonIntro struct {
	ID          uint    `json:"id"`
	ProviderID  uint    `json:"-" gorm:"column:provider_id;not null;unique_index:idx_provider_person"`
	Name        string  `json:"name"`
	Gender      uint8   `json:"gender"`
	Order       *uint32 `json:"order,omitempty" gorm:"-"`
	ProfilePath string  `json:"profilePath"`
}

func (PersonIntro) TableName() string {
	return "people"
}
