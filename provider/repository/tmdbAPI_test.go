package repository

import (
	"reflect"
	"testing"

	"github.com/kunhou/TMDB/config"
	"github.com/kunhou/TMDB/models"
)

var cfg = config.GetConfig()

func Test_tmdbRepository_GetMovieLastID(t *testing.T) {
	type fields struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			"ok",
			fields{cfg.TMDBToken},
			537141,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdb := &tmdbRepository{
				token: tt.fields.token,
			}
			got, err := tmdb.GetMovieLastID()
			if (err != nil) != tt.wantErr {
				t.Errorf("tmdbRepository.GetMovieLastID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tmdbRepository.GetMovieLastID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tmdbRepository_GetMovieDetail(t *testing.T) {
	type fields struct {
		token string
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Movie
		wantErr bool
	}{
		{
			"ok",
			fields{cfg.TMDBToken},
			args{537141},
			&models.Movie{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdb := &tmdbRepository{
				token: tt.fields.token,
			}
			got, err := tmdb.GetMovieDetail(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("tmdbRepository.GetMovieDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tmdbRepository.GetMovieDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tmdbRepository_GetPersonLastID(t *testing.T) {
	type fields struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			"ok",
			fields{cfg.TMDBToken},
			537141,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdb := &tmdbRepository{
				token: tt.fields.token,
			}
			got, err := tmdb.GetPersonLastID()
			if (err != nil) != tt.wantErr {
				t.Errorf("tmdbRepository.GetPersonLastID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tmdbRepository.GetPersonLastID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tmdbRepository_GetPersonDetail(t *testing.T) {
	type fields struct {
		token string
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Person
		wantErr bool
	}{
		{
			"ok",
			fields{cfg.TMDBToken},
			args{1},
			&models.Person{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmdb := &tmdbRepository{
				token: tt.fields.token,
			}
			got, err := tmdb.GetPersonDetail(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("tmdbRepository.GetPersonDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tmdbRepository.GetPersonDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
