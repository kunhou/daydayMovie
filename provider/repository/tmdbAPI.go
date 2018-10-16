package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/provider"
	"github.com/pkg/errors"
)

var (
	apiURL            = "https://api.themoviedb.org/3/"
	DISCOVER_PATH     = "/discover/movie"
	LATEST_MOVIE_ID   = "/movie/latest"
	GET_MOVIE_DETAIL  = "/movie/%d"
	LATEST_PERSON_ID  = "/person/latest"
	GET_PERSON_DETAIL = "/person/%d"
	LATEST_TV_ID      = "/tv/latest"
	GET_TV_DETAIL     = "/tv/%d"
	GET_TV_SEASON     = "/tv/%d/season/%d"
	GET_MOVIE_CREDITs = "/movie/%d/credits"
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

func (tmdb *tmdbRepository) GetPersonLastID() (int, error) {
	type responseBody struct {
		ID int `json:"id"`
	}
	var data responseBody
	if err := tmdb.request(LATEST_PERSON_ID, nil, &data); err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (tmdb *tmdbRepository) GetPersonDetail(id int) (*models.Person, error) {
	type responseBody struct {
		ID                 uint     `json:"id"`
		Birthday           string   `json:"birthday"`
		KnownForDepartment string   `json:"known_for_department"`
		Deathday           string   `json:"deathday"`
		Name               string   `json:"name"`
		AlsoKnownAs        []string `json:"also_known_as"`
		Gender             uint8    `json:"gender"`
		Biography          string   `json:"biography"`
		Popularity         float32  `json:"popularity"`
		PlaceOfBirth       string   `json:"place_of_birth"`
		ProfilePath        string   `json:"profile_path"`
		Adult              bool     `json:"adult"`
		ImdbID             string   `json:"imdb_id"`
		Homepage           string   `json:"homepage"`
	}

	var data responseBody
	urlPath := fmt.Sprintf(GET_PERSON_DETAIL, id)
	if err := tmdb.request(urlPath, nil, &data); err != nil {
		return nil, err
	}

	var birthday, deathday *time.Time
	birthday = nil
	if len(data.Birthday) > 0 {
		t, err := time.Parse("2006-01-02", data.Birthday)
		if err != nil {
			log.WithError(err).Error("Parse time error : " + data.Birthday)
		} else {
			birthday = &t
		}
	}
	deathday = nil
	if len(data.Deathday) > 0 {
		t, err := time.Parse("2006-01-02", data.Deathday)
		if err != nil {
			log.WithError(err).Error("Parse time error : " + data.Deathday)
		} else {
			deathday = &t
		}
	}

	return &models.Person{
		Provider:           "tmdb",
		ProviderID:         data.ID,
		Name:               data.Name,
		Birthday:           birthday,
		Deathday:           deathday,
		Gender:             data.Gender,
		Biography:          data.Biography,
		Popularity:         data.Popularity,
		PlaceOfBirth:       data.PlaceOfBirth,
		Adult:              data.Adult,
		ImdbID:             data.ImdbID,
		Homepage:           data.Homepage,
		AlsoKnownAs:        data.AlsoKnownAs,
		ProfilePath:        data.ProfilePath,
		KnownForDepartment: data.KnownForDepartment,
	}, nil
}

func (tmdb *tmdbRepository) GetTVLastID() (int, error) {
	type responseBody struct {
		ID int `json:"id"`
	}
	var data responseBody
	if err := tmdb.request(LATEST_TV_ID, nil, &data); err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (tmdb *tmdbRepository) GetTVDetail(id int) (*models.TV, error) {
	type responseBody struct {
		BackdropPath   string             `json:"backdrop_path"`
		CreatedBy      []models.CreatedBy `json:"created_by"`
		EpisodeRunTime []int64            `json:"episode_run_time"`
		FirstAirDate   string             `json:"first_air_date"`
		Genres         []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"genres"`
		Homepage         string   `json:"homepage"`
		ID               uint     `json:"id"`
		InProduction     bool     `json:"in_production"`
		Languages        []string `json:"languages"`
		LastAirDate      string   `json:"last_air_date"`
		LastEpisodeToAir struct {
			AirDate        string  `json:"air_date"`
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
		} `json:"last_episode_to_air"`
		Name             string `json:"name"`
		NextEpisodeToAir struct {
			AirDate        string  `json:"air_date"`
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
		} `json:"next_episode_to_air"`
		Networks []struct {
			Name          string `json:"name"`
			ID            int    `json:"id"`
			LogoPath      string `json:"logo_path"`
			OriginCountry string `json:"origin_country"`
		} `json:"networks"`
		NumberOfEpisodes    int      `json:"number_of_episodes"`
		NumberOfSeasons     int      `json:"number_of_seasons"`
		OriginCountry       []string `json:"origin_country"`
		OriginalLanguage    string   `json:"original_language"`
		OriginalName        string   `json:"original_name"`
		Overview            string   `json:"overview"`
		Popularity          float64  `json:"popularity"`
		PosterPath          string   `json:"poster_path"`
		ProductionCompanies []struct {
			ID            int    `json:"id"`
			LogoPath      string `json:"logo_path"`
			Name          string `json:"name"`
			OriginCountry string `json:"origin_country"`
		} `json:"production_companies"`
		Seasons []struct {
			AirDate      string `json:"air_date"`
			EpisodeCount int    `json:"episode_count"`
			ID           int    `json:"id"`
			Name         string `json:"name"`
			Overview     string `json:"overview"`
			PosterPath   string `json:"poster_path"`
			SeasonNumber int    `json:"season_number"`
		} `json:"seasons"`
		Status      string  `json:"status"`
		Type        string  `json:"type"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
	}

	var data responseBody
	urlPath := fmt.Sprintf(GET_TV_DETAIL, id)
	if err := tmdb.request(urlPath, nil, &data); err != nil {
		return nil, err
	}
	var firstAirDate, lastAirDate *time.Time
	firstAirDate, lastAirDate = nil, nil
	if len(data.FirstAirDate) > 0 {
		t, err := time.Parse("2006-01-02", data.FirstAirDate)
		if err != nil {
			log.WithError(err).Error("Parse time error : " + data.FirstAirDate)
		} else {
			firstAirDate = &t
		}
	}
	if len(data.LastAirDate) > 0 {
		t, err := time.Parse("2006-01-02", data.LastAirDate)
		if err != nil {
			log.WithError(err).Error("Parse time error : " + data.LastAirDate)
		} else {
			lastAirDate = &t
		}
	}

	genreIDs := []int64{}
	for _, g := range data.Genres {
		genreIDs = append(genreIDs, g.ID)
	}
	seasons := []models.Season{}
	for _, s := range data.Seasons {
		var airDate *time.Time
		airDate = nil
		if len(s.AirDate) > 0 {
			t, err := time.Parse("2006-01-02", s.AirDate)
			if err != nil {
				log.WithError(err).Error("Parse time error : " + s.AirDate)
			} else {
				airDate = &t
			}
		}
		seasons = append(seasons, models.Season{
			AirDate:      airDate,
			EpisodeCount: s.EpisodeCount,
			Name:         s.Name,
			Overview:     s.Overview,
			PosterPath:   s.PosterPath,
			SeasonNumber: s.SeasonNumber,
		})
	}
	t := models.TV{
		Provider:       "tmdb",
		ProviderID:     data.ID,
		BackdropPath:   data.BackdropPath,
		CreatedBy:      data.CreatedBy,
		EpisodeRunTime: data.EpisodeRunTime,
		FirstAirDate:   firstAirDate,
		GenreIds:       genreIDs,
		Homepage:       data.Homepage,
		InProduction:   data.InProduction,
		LastAirDate:    lastAirDate,
		LastEpisodeToAir: models.LastEpisodeToAir{
			AirDate:        data.LastEpisodeToAir.AirDate,
			EpisodeNumber:  data.LastEpisodeToAir.EpisodeNumber,
			ID:             data.LastEpisodeToAir.ID,
			Name:           data.LastEpisodeToAir.Name,
			Overview:       data.LastEpisodeToAir.Overview,
			ProductionCode: data.LastEpisodeToAir.ProductionCode,
			SeasonNumber:   data.LastEpisodeToAir.SeasonNumber,
			ShowID:         data.LastEpisodeToAir.ShowID,
			StillPath:      data.LastEpisodeToAir.StillPath,
			VoteAverage:    data.LastEpisodeToAir.VoteAverage,
			VoteCount:      data.LastEpisodeToAir.VoteCount,
		},
		Name: data.Name,
		NextEpisodeToAir: models.NextEpisodeToAir{
			ID:             data.NextEpisodeToAir.ID,
			AirDate:        data.NextEpisodeToAir.AirDate,
			EpisodeNumber:  data.NextEpisodeToAir.EpisodeNumber,
			Name:           data.NextEpisodeToAir.Name,
			Overview:       data.NextEpisodeToAir.Overview,
			ProductionCode: data.NextEpisodeToAir.ProductionCode,
			SeasonNumber:   data.NextEpisodeToAir.SeasonNumber,
			ShowID:         data.NextEpisodeToAir.ShowID,
			StillPath:      data.NextEpisodeToAir.StillPath,
			VoteAverage:    data.NextEpisodeToAir.VoteAverage,
			VoteCount:      data.NextEpisodeToAir.VoteCount,
		},
		NumberOfEpisodes: data.NumberOfEpisodes,
		NumberOfSeasons:  data.NumberOfSeasons,
		OriginalLanguage: data.OriginalLanguage,
		OriginalName:     data.OriginalName,
		Overview:         data.Overview,
		Popularity:       data.Popularity,
		PosterPath:       data.PosterPath,
		Status:           data.Status,
		Type:             data.Type,
		VoteAverage:      data.VoteAverage,
		VoteCount:        data.VoteCount,
		OriginCountry:    data.OriginCountry,
		Seasons:          seasons,
	}

	return &t, nil
}

func (tmdb *tmdbRepository) GetTVSeasonVote(tvID uint, seasonID int) (voteAvg float64, voteCount int, err error) {
	type body struct {
		Episodes []struct {
			AirDate        string        `json:"air_date"`
			EpisodeNumber  int           `json:"episode_number"`
			ID             int           `json:"id"`
			Name           string        `json:"name"`
			Overview       string        `json:"overview"`
			ProductionCode interface{}   `json:"production_code"`
			SeasonNumber   int           `json:"season_number"`
			ShowID         int           `json:"show_id"`
			StillPath      interface{}   `json:"still_path"`
			VoteAverage    float64       `json:"vote_average"`
			VoteCount      int           `json:"vote_count"`
			Crew           []interface{} `json:"crew"`
			GuestStars     []interface{} `json:"guest_stars"`
		} `json:"episodes"`
	}
	var data body
	urlPath := fmt.Sprintf(GET_TV_SEASON, tvID, seasonID)
	if err := tmdb.request(urlPath, nil, &data); err != nil {
		return 0, 0, err
	}
	voteTotalPoint := float64(0)
	for _, e := range data.Episodes {
		voteCount += e.VoteCount
		voteTotalPoint += e.VoteAverage * float64(e.VoteCount)
	}
	if voteCount > 0 {
		voteAvg = math.Round(voteTotalPoint*100/float64(voteCount)) / float64(100)
	}
	return voteAvg, voteCount, nil
}

func (tmdb *tmdbRepository) GetMovieCredits(movieID uint) (casts []models.Credit, crews []models.Credit, err error) {
	casts = []models.Credit{}
	crews = []models.Credit{}
	type body struct {
		ID   int `json:"id"`
		Cast []struct {
			CastID      int    `json:"cast_id"`
			Character   string `json:"character"`
			CreditID    string `json:"credit_id"`
			Gender      int    `json:"gender"`
			ID          uint   `json:"id"`
			Name        string `json:"name"`
			Order       int    `json:"order"`
			ProfilePath string `json:"profile_path"`
		} `json:"cast"`
		Crew []struct {
			CreditID    string `json:"credit_id"`
			Department  string `json:"department"`
			Gender      int    `json:"gender"`
			ID          uint   `json:"id"`
			Job         string `json:"job"`
			Name        string `json:"name"`
			ProfilePath string `json:"profile_path"`
		} `json:"crew"`
	}
	var data body
	urlPath := fmt.Sprintf(GET_MOVIE_CREDITs, movieID)
	if err := tmdb.request(urlPath, nil, &data); err != nil {
		return casts, crews, err
	}
	for i := range data.Cast {
		casts = append(casts, models.Credit{
			ProviderPersonID: data.Cast[i].ID,
			Order:            &data.Cast[i].Order,
			Character:        &data.Cast[i].Character,
			Type:             models.CreditTypeCast,
			Cast:             models.CastMovie,
			CastID:           movieID,
		})
	}
	for i := range data.Crew {
		crews = append(crews, models.Credit{
			ProviderPersonID: data.Crew[i].ID,
			Department:       &data.Crew[i].Department,
			Type:             models.CreditTypeCrew,
			Cast:             models.CastMovie,
			CastID:           movieID,
		})
	}
	return casts, crews, nil
}

var requestMutex sync.Mutex

const CrawlerInterval = 700 * time.Millisecond

func (tmdb *tmdbRepository) request(urlPath string, options map[string]string, v interface{}) error {
	time.Sleep(CrawlerInterval)
	requestMutex.Lock()
	defer requestMutex.Unlock()
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
