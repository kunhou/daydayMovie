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

func Test_pgsqlRepository_PeopleList(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		page   int
		limit  int
		order  map[string]string
		search map[string]interface{}
	}
	btime := time.Date(2018, time.June, 24, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		p       pgsqlRepository
		args    args
		want    []*models.Person
		want1   *models.Page
		wantErr bool
	}{
		{
			"ok",
			pgsqlRepository{db.DB},
			args{2, 10, map[string]string{"popularity": "desc"}, map[string]interface{}{
				"birthday": btime,
			}},
			nil,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.p.PeopleList(tt.args.page, tt.args.limit, tt.args.order, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.PeopleList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgsqlRepository.PeopleList() got = %#v, want %v", got, tt.want)
				for _, m := range got {
					t.Errorf(" movie = %+v", m)
				}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("pgsqlRepository.PeopleList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_pgsqlRepository_CreditStore(t *testing.T) {
	type args struct {
		c *models.Credit
	}
	o1 := 51
	o2 := 52
	d1 := "directing"
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
				&models.Credit{
					PersonID:   1211,
					Cast:       "tv",
					Type:       "crew",
					CastID:     12,
					Department: &d1,
				},
			},
			0,
			false,
		},
		{
			"ok",
			pgsqlRepository{db.DB},
			args{
				&models.Credit{
					PersonID: 1210,
					Cast:     "movie",
					Type:     "cast",
					CastID:   12,
					Order:    &o1,
				},
			},
			0,
			false,
		},
		{
			"ok",
			pgsqlRepository{db.DB},
			args{
				&models.Credit{
					PersonID: 1210,
					Cast:     "movie",
					Type:     "cast",
					CastID:   12,
					Order:    &o2,
				},
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.CreditStore(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgsqlRepository.CreditStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pgsqlRepository.CreditStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
