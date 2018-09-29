package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kunhou/TMDB/httputil"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/movie"
)

type HttpMovieHandler struct {
	MUsecase movie.MovieUsecase
}

type movieListResponse struct {
	*models.Page
	Results []*models.MovieIntro `json:"results"`
}
type tvListResponse struct {
	*models.Page
	Results []*models.TVIntro `json:"results"`
}

type peopleListResponse struct {
	*models.Page
	Results []*models.Person `json:"results"`
}

func NewMovieHttpHandler(mu movie.MovieUsecase) *HttpMovieHandler {
	handler := &HttpMovieHandler{
		MUsecase: mu,
	}
	return handler
}

func (ph *HttpMovieHandler) MovieList(c *gin.Context) {
	var page, limit int
	orderBy := make(map[string]string)
	if p, ok := c.GetQuery("page"); ok {
		if pageInt, err := strconv.Atoi(p); err == nil {
			page = pageInt
		}
	}
	if l, ok := c.GetQuery("limit"); ok {
		if limitInt, err := strconv.Atoi(l); err == nil {
			limit = limitInt
		}
	}
	if sb, ok := c.GetQuery("sort_by"); ok {
		sbs := strings.Split(sb, ".")
		if len(sbs) == 2 {
			orderBy[sbs[0]] = sbs[1]
		}
	}
	movieList, pageInfo, err := ph.MUsecase.MovieList(page, limit, orderBy)
	if err != nil {
		httputil.ResponseFail(c, http.StatusInternalServerError, 4001, "Internal Server error while fetching movie list", err)
		return
	}

	c.JSON(http.StatusOK, movieListResponse{
		pageInfo,
		movieList,
	})
}

func (ph *HttpMovieHandler) MovieDetail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		httputil.ResponseFail(c, http.StatusNotFound, 4002, "Invalid Path", err)
		return
	}
	movie, err := ph.MUsecase.MovieDetail(uint(idInt))
	if err != nil {
		httputil.ResponseFail(c, http.StatusInternalServerError, 4001, "Internal Server error while fetching movie detail", err)
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (ph *HttpMovieHandler) TVList(c *gin.Context) {
	var page, limit int
	orderBy := make(map[string]string)
	if p, ok := c.GetQuery("page"); ok {
		if pageInt, err := strconv.Atoi(p); err == nil {
			page = pageInt
		}
	}
	if l, ok := c.GetQuery("limit"); ok {
		if limitInt, err := strconv.Atoi(l); err == nil {
			limit = limitInt
		}
	}
	if sb, ok := c.GetQuery("sort_by"); ok {
		sbs := strings.Split(sb, ".")
		if len(sbs) == 2 {
			orderBy[sbs[0]] = sbs[1]
		}
	}
	tvList, pageInfo, err := ph.MUsecase.TVList(page, limit, orderBy)
	if err != nil {
		httputil.ResponseFail(c, http.StatusInternalServerError, 4001, "Internal Server error while fetching movie list", err)
		return
	}

	c.JSON(http.StatusOK, tvListResponse{
		pageInfo,
		tvList,
	})
}

func (ph *HttpMovieHandler) PeopleList(c *gin.Context) {
	var page, limit int
	orderBy := make(map[string]string)
	if p, ok := c.GetQuery("page"); ok {
		if pageInt, err := strconv.Atoi(p); err == nil {
			page = pageInt
		}
	}
	if l, ok := c.GetQuery("limit"); ok {
		if limitInt, err := strconv.Atoi(l); err == nil {
			limit = limitInt
		}
	}
	if sb, ok := c.GetQuery("sort_by"); ok {
		sbs := strings.Split(sb, ".")
		if len(sbs) == 2 {
			orderBy[sbs[0]] = sbs[1]
		}
	}
	search := map[string]interface{}{}
	if birthday, ok := c.GetQuery("birthday"); ok {
		tBirthday, err := time.Parse("01-02", birthday)
		if err != nil {
			httputil.ResponseFail(c, http.StatusBadRequest, 3001, "Birthday format error", err)
			return
		}
		search["birthday"] = tBirthday
	}

	peopleList, pageInfo, err := ph.MUsecase.PeopleList(page, limit, orderBy, search)
	if err != nil {
		httputil.ResponseFail(c, http.StatusInternalServerError, 4001, "Internal Server error while fetching movie list", err)
		return
	}

	c.JSON(http.StatusOK, peopleListResponse{
		pageInfo,
		peopleList,
	})
}
