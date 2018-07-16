package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/kunhou/TMDB/config"
	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/provider"
	"github.com/pkg/errors"
)

var (
	cfg              = config.GetConfig()
	apiURL           = "https://api.themoviedb.org/3/"
	DISCOVER_PATH    = "/discover/movie"
	LATEST_MOVIE_ID  = "/movie/latest"
	GET_MOVIE_DETAIL = "/movie/%d"
)

type tmdbRepository struct {
	token string
}

func NewTMDBRepository(token string) provider.ProviderRepository {
	return &tmdbRepository{token}
}

func (tmdb *tmdbRepository) GetMovieTotalPages() (int, error) {
	type responseBody struct {
		Page         int `json:"page"`
		TotalResults int `json:"total_results"`
		TotalPages   int `json:"total_pages"`
	}
	var data responseBody
	var options = make(map[string]string)
	options["page"] = "1"
	if err := tmdb.request(DISCOVER_PATH, options, &data); err != nil {
		return 0, err
	}
	return data.TotalPages, nil
}

func (tmdb *tmdbRepository) GetMovieLastID() (int, error) {
	type responseBody struct {
		ID int `json:"id"`
	}
	var data responseBody
	if err := tmdb.request(LATEST_MOVIE_ID, nil, &data); err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (tmdb *tmdbRepository) GetMovieDetail(id int) (*models.Movie, error) {
	type responseBody struct {
		ID               uint    `json:"id"`
		VoteCount        int     `json:"vote_count"`
		Video            bool    `json:"video"`
		VoteAverage      float32 `json:"vote_average"`
		Title            string  `json:"title"`
		Popularity       float32 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		GenreIds         []struct {
			ID   int64  `json:"id"`
			Name string `json:"name`
		} `json:"genres"`
		BackdropPath string `json:"backdrop_path"`
		Adult        bool   `json:"adult"`
		Overview     string `json:"overview"`
		ReleaseDate  string `json:"release_date"`
	}

	var data responseBody
	urlPath := fmt.Sprintf(GET_MOVIE_DETAIL, id)
	if err := tmdb.request(urlPath, nil, &data); err != nil {
		return nil, err
	}

	var rTime *time.Time
	rTime = nil
	if len(data.ReleaseDate) > 0 {
		t, err := time.Parse("2006-01-02", data.ReleaseDate)
		if err != nil {
			log.WithError(err).Error("Parse time error : " + data.ReleaseDate)
		} else {
			rTime = &t
		}
	}
	genreIDs := []int64{}
	for _, g := range data.GenreIds {
		genreIDs = append(genreIDs, g.ID)
	}

	return &models.Movie{
		Provider:         "tmdb",
		ProviderID:       data.ID,
		Title:            data.Title,
		OriginalTitle:    data.OriginalTitle,
		Popularity:       data.Popularity,
		VoteAverage:      data.VoteAverage,
		VoteCount:        data.VoteCount,
		PosterPath:       data.PosterPath,
		OriginalLanguage: data.OriginalLanguage,
		GenreIds:         genreIDs,
		BackdropPath:     data.BackdropPath,
		Adult:            data.Adult,
		Overview:         data.Overview,
		ReleaseDate:      rTime,
	}, nil
}

func (tmdb *tmdbRepository) GetMovieWithPage(page int) ([]*models.Movie, error) {
	type responseBody struct {
		Page         int `json:"page"`
		TotalResults int `json:"total_results"`
		TotalPages   int `json:"total_pages"`
		Results      []struct {
			ID               uint    `json:"id"`
			VoteCount        int     `json:"vote_count"`
			Video            bool    `json:"video"`
			VoteAverage      float32 `json:"vote_average"`
			Title            string  `json:"title"`
			Popularity       float32 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			GenreIds         []int64 `json:"genre_ids"`
			BackdropPath     string  `json:"backdrop_path"`
			Adult            bool    `json:"adult"`
			Overview         string  `json:"overview"`
			ReleaseDate      string  `json:"release_date"`
		} `json:"results"`
	}
	if page < 1 {
		page = 1
	}
	var data responseBody
	var options = make(map[string]string)
	var movies []*models.Movie
	options["page"] = strconv.Itoa(page)
	if err := tmdb.request(DISCOVER_PATH, options, &data); err != nil {
		return movies, err
	}
	for _, movie := range data.Results {
		var rTime *time.Time
		rTime = nil
		if len(movie.ReleaseDate) > 0 {
			t, err := time.Parse("2006-01-02", movie.ReleaseDate)
			if err != nil {
				log.WithError(err).Error("Parse time error : " + movie.ReleaseDate)
			} else {
				rTime = &t
			}
		}
		movies = append(movies, &models.Movie{
			Provider:         "tmdb",
			ProviderID:       movie.ID,
			Title:            movie.Title,
			OriginalTitle:    movie.OriginalTitle,
			Popularity:       movie.Popularity,
			VoteAverage:      movie.VoteAverage,
			VoteCount:        movie.VoteCount,
			PosterPath:       movie.PosterPath,
			OriginalLanguage: movie.OriginalLanguage,
			GenreIds:         movie.GenreIds,
			BackdropPath:     movie.BackdropPath,
			Adult:            movie.Adult,
			Overview:         movie.Overview,
			ReleaseDate:      rTime,
		})
	}
	return movies, nil
}

func (tmdb *tmdbRepository) request(urlPath string, options map[string]string, v interface{}) error {
	u, err := url.Parse(apiURL)
	u.Path = path.Join(u.Path, urlPath)
	q := url.Values{}
	q.Add("api_key", tmdb.token)
	q.Add("language", "zh-TW")
	// q.Add("region", "TW")
	for k, v := range options {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return errors.Wrap(err, "http get error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "ReadAll error")
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return provider.APINotFoundError{
				Path: urlPath,
			}
		}
		return errors.New("fail: " + string(body))
	}
	if err := json.Unmarshal(body, &v); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url":    resp.Request.Host,
			"status": resp.Status,
			"body":   string(body),
		}).Error("Unmarshal Fail")
		return errors.Wrap(err, "Unmarshal error")
	}

	return nil
}
