package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kunhou/TMDB/db"
	"github.com/kunhou/TMDB/models"
)

func Test_pgsqlRepository_Store(t *testing.T) {
	type args struct {
		m *models.Movie
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
				&models.Movie{
					ProviderID:       1111,
					Provider:         "tmdb",
					Title:            "1233333",
					OriginalTitle:    "bb",
					Popularity:       2.313,
					VoteAverage:      2.22,
					VoteCount:        222,
					PosterPath:       "/ssss/.ss",
					OriginalLanguage: "en",
					GenreIds:         []int64{1, 2, 3},
					BackdropPath:     "aaa",
					Adult:            true,
					Overview:         "aaa",
					ReleaseDate:      &timeNow,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Store(tt.args.m)
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
		movies []*models.Movie
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
				[]*models.Movie{
					&models.Movie{
						ProviderID:       11,
						Provider:         "tmdb",
						Title:            "11",
						OriginalTitle:    "bb",
						Popularity:       2.313,
						VoteAverage:      2.22,
						VoteCount:        222,
						PosterPath:       "/aaa/.ss",
						OriginalLanguage: "en",
						GenreIds:         []int64{1, 2, 3, 4, 5},
						BackdropPath:     "aaa",
						Adult:            true,
						Overview:         "aaa",
						ReleaseDate:      &timeNow,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
					&models.Movie{
						ProviderID:       12,
						Provider:         "tmdb",
						Title:            "12",
						OriginalTitle:    "bbbb",
						Popularity:       2.313,
						VoteAverage:      2.22,
						VoteCount:        222,
						PosterPath:       "/aaa/.ss",
						OriginalLanguage: "en",
						GenreIds:         []int64{1, 2, 3, 4, 5},
						BackdropPath:     "aaa",
						Adult:            true,
						Overview:         "aaa",
						ReleaseDate:      &timeNow,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
					&models.Movie{
						ProviderID:       13,
						Provider:         "tmdb",
						Title:            "13",
						OriginalTitle:    "bbbb",
						Popularity:       2.313,
						VoteAverage:      2.22,
						VoteCount:        222,
						PosterPath:       "/aaa/.ss",
						OriginalLanguage: "en",
						GenreIds:         []int64{1, 2, 3, 4, 5},
						BackdropPath:     "aaa",
						Adult:            true,
						Overview:         "aaa",
						ReleaseDate:      &timeNow,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
					&models.Movie{
						ProviderID:       14,
						Provider:         "tmdb",
						Title:            "14",
						OriginalTitle:    "ccbccc",
						Popularity:       2.313,
						VoteAverage:      2.22,
						VoteCount:        222,
						PosterPath:       "/aaa/.ss",
						OriginalLanguage: "en",
						GenreIds:         []int64{1, 2, 3, 4, 5},
						BackdropPath:     "aaa",
						Adult:            true,
						Overview:         "aaa",
						ReleaseDate:      &timeNow,
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.BatchStore(tt.args.movies); (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.BatchStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgsqlRepository_List(t *testing.T) {
	type args struct {
		page  int
		limit int
		query map[string]string
	}
	tests := []struct {
		name    string
		p       pgsqlRepository
		args    args
		want    []*models.MovieIntro
		want2   *models.Page
		wantErr bool
	}{
		{
			"ok",
			pgsqlRepository{db.DB},
			args{2, 10, map[string]string{"popularity": "desc"}},
			nil,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, err := tt.p.MovieList(tt.args.page, tt.args.limit, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.MovieList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgsqlRepository.MovieList() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("pgsqlRepository.MovieList() = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}

func Test_pgsqlRepository_MovieDetail(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		p       pgsqlRepository
		args    args
		want    *models.Movie
		wantErr bool
	}{
		{
			"ok",
			pgsqlRepository{db.DB},
			args{1},
			nil,
			false,
		},
		{
			"ok",
			pgsqlRepository{db.DB},
			args{10},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MovieDetail(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.MovieDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgsqlRepository.MovieDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgsqlRepository_TVStore(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		t *models.TV
	}
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
				&models.TV{
					Provider:     "tmdb",
					ProviderID:   1121,
					BackdropPath: "asdasd",
					CreatedBy: []models.CreatedBy{
						models.CreatedBy{123, "asd", "aaaa", 1, "wdds"},
					},
					EpisodeRunTime: []int64{54},
					GenreIds:       []int64{54},
					OriginCountry:  []string{},
					VoteAverage:    12.2,
					VoteCount:      12,
				},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.TVStore(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.TVStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pgsqlRepository.TVStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
