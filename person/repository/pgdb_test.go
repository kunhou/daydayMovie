package repository

import (
	"testing"
	"time"

	"github.com/kunhou/TMDB/db"
	"github.com/kunhou/TMDB/models"
)

func Test_pgsqlRepository_Store(t *testing.T) {
	type args struct {
		pn *models.Person
	}
	timeNow := time.Now()
	tests := []struct {
		name    string
		p       pgsqlRepository
		args    args
		want    uint
		wantErr bool
	}{
		{
			"ok",
			pgsqlRepository{db.DB},
			args{
				&models.Person{
					Name:         "test234",
					ProviderID:   2222,
					Provider:     "tmdb",
					Birthday:     &timeNow,
					Deathday:     &timeNow,
					Gender:       2,
					Biography:    "abcdef",
					Popularity:   2.333,
					PlaceOfBirth: "taiwan",
					Adult:        true,
					ImdbID:       "imdbbb",
					Homepage:     "homepage",
					AlsoKnownAs:  []string{"1111"},
					ProfilePath:  "aaaaa",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Store(tt.args.pn)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pgsqlRepository.Store() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgsqlRepository_BatchStore(t *testing.T) {
	type args struct {
		people []*models.Person
	}
	timeNow := time.Now()
	tests := []struct {
		name    string
		p       pgsqlRepository
		args    args
		wantErr bool
	}{
		{
			"ok",
			pgsqlRepository{db.DB},
			args{
				[]*models.Person{
					&models.Person{
						Name:         "test234",
						ProviderID:   6666,
						Provider:     "tmdb",
						Birthday:     &timeNow,
						Deathday:     &timeNow,
						Gender:       2,
						Biography:    "abcdef",
						Popularity:   2.333,
						PlaceOfBirth: "taiwan",
						Adult:        true,
						ImdbID:       "imdbbb",
						Homepage:     "homepage",
						AlsoKnownAs:  []string{"1111"},
						ProfilePath:  "aaaaa",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
					&models.Person{
						Name:         "test22224",
						ProviderID:   7777,
						Provider:     "tmdb",
						Birthday:     &timeNow,
						Deathday:     &timeNow,
						Gender:       2,
						Biography:    "abcdef",
						Popularity:   2.333,
						PlaceOfBirth: "taiwan",
						Adult:        true,
						ImdbID:       "imdbbb",
						Homepage:     "homepage",
						AlsoKnownAs:  []string{"1111"},
						ProfilePath:  "aaaaa",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.BatchStore(tt.args.people); (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.BatchStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
