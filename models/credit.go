package models

const (
	CastTV              = "tv"
	CastMovie           = "movie"
	CreditTypeCast      = "cast"
	CreditTypeCrew      = "crew"
	DepartmentDirecting = "Directing"
)

type Credit struct {
	ID               uint   `gorm:"primary_key"`
	PersonID         uint   `gorm:"not null;unique_index:idx_person_cast_type"`
	ProviderPersonID uint   `gorm:"-"`
	CastID           uint   `gorm:"not null;unique_index:idx_person_cast_type"`
	Cast             string `gorm:"type:varchar(255);not null;unique_index:idx_person_cast_type"`
	Type             string `gorm:"type:varchar(255);not null;unique_index:idx_person_cast_type"`
	Order            *int
	Character        *string `gorm:"type:varchar(255)"`
	Department       *string `gorm:"type:varchar(255)"`
}
