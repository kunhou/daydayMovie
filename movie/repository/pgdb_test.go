package repository

import (
	"testing"
	"time"

	"github.com/kunhou/TMDB/db"
	"github.com/kunhou/TMDB/models"
)

func Test_pgsqlRepository_Store(t *testing.T) {
	type args struct {
		m *models.Movie
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
					ReleaseDate:      time.Now(),
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
		movies []models.Movie
	}
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
				[]models.Movie{
					models.Movie{
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
						ReleaseDate:      time.Now(),
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
					models.Movie{
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
						ReleaseDate:      time.Now(),
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
					models.Movie{
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
						ReleaseDate:      time.Now(),
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					},
					models.Movie{
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
						ReleaseDate:      time.Now(),
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
