package usecase

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/kunhou/TMDB/log"
	"github.com/kunhou/TMDB/models"
	"github.com/kunhou/TMDB/movie"
	"github.com/kunhou/TMDB/person"
	"github.com/kunhou/TMDB/provider"
)

const CrawlerInterval = 700 * time.Millisecond

type tmdbUsecase struct {
	providerRepo provider.ProviderRepository
	movieRepo    movie.MovieRepository
	personRepo   person.PersonRepository
}

func NewTmdbUsecase(repo provider.ProviderRepository, movieRepo movie.MovieRepository, personRepo person.PersonRepository) provider.ProviderUsecase {
	return &tmdbUsecase{repo, movieRepo, personRepo}
}

func (tu *tmdbUsecase) StartCrawlerMovie(ch chan *models.Movie) {
	lastestID, err := tu.providerRepo.GetMovieLastID()
	if err != nil {
		log.WithError(err).Error("Get LastID Fail")
	}
	for id := lastestID; id > 0; id-- {
		// time.Sleep(CrawlerInterval)
		m, err := tu.providerRepo.GetMovieDetail(id)
		if err != nil {
			if _, ok := err.(provider.APINotFoundError); !ok {
				log.WithError(err).Error("Get discover Fail")
			}
			continue
		}
		ch <- m
	}
	return
}

func (tu *tmdbUsecase) StartCrawlerPerson(ch chan *models.Person) {
	lastestID, err := tu.providerRepo.GetPersonLastID()
	if err != nil {
		log.WithError(err).Error("Get LastID Fail")
	}
	for id := lastestID; id > 0; id-- {
		// time.Sleep(CrawlerInterval)
		person, err := tu.providerRepo.GetPersonDetail(id)
		if err != nil {
			if _, ok := err.(provider.APINotFoundError); !ok {
				log.WithError(err).Error("Get discover Fail")
			}
			continue
		}
		ch <- person
	}
	return
}

var movieWriteChan chan *models.Movie
var movieWriteSyncOnce sync.Once

func (tu *tmdbUsecase) CreateBatchStoreMovieTask() chan *models.Movie {
	movieWriteSyncOnce.Do(func() {
		movieWriteChan = make(chan *models.Movie, 100000)
		movieBuff := []*models.Movie{}
		go func() {
			t := time.NewTicker(1 * time.Minute)
			for {
				select {
				case <-t.C:
					if len(movieBuff) == 0 {
						continue
					}
					if err := tu.movieRepo.BatchStore(movieBuff); err != nil {
						log.WithError(err).Error("Movie Task store fail")
					}
					movieBuff = []*models.Movie{}
				case movie := <-movieWriteChan:
					movieBuff = append(movieBuff, movie)
					if len(movieBuff) > 100 {
						if err := tu.movieRepo.BatchStore(movieBuff); err != nil {
							log.WithError(err).Error("Movie Task store fail")
						}
						movieBuff = []*models.Movie{}
					}
				}
			}
		}()
	})

	return movieWriteChan
}

var personWriteChan chan *models.Person
var personWriteSyncOnce sync.Once

func (tu *tmdbUsecase) CreateBatchStorePersonTask() chan *models.Person {
	personWriteSyncOnce.Do(func() {
		personWriteChan = make(chan *models.Person, 100000)
		personBuff := []*models.Person{}
		go func() {
			t := time.NewTicker(1 * time.Minute)
			for {
				select {
				case <-t.C:
					if len(personBuff) == 0 {
						continue
					}
					if err := tu.personRepo.BatchStore(personBuff); err != nil {
						log.WithError(err).Error("person Task store fail")
					}
					personBuff = []*models.Person{}
				case person := <-personWriteChan:
					personBuff = append(personBuff, person)
					if len(personBuff) > 1 {
						if err := tu.personRepo.BatchStore(personBuff); err != nil {
							log.WithError(err).Error("person Task store fail")
						}
						personBuff = []*models.Person{}
					}
				}
			}
		}()
	})

	return personWriteChan
}

func (tu *tmdbUsecase) StartCrawlerTV(ch chan *models.TV) {
	lastestID, err := tu.providerRepo.GetTVLastID()
	log.Debug("tv lastestID: ", lastestID)
	if err != nil {
		log.WithError(err).Error("Get LastID Fail")
	}
	for id := lastestID; id > 0; id-- {
		// time.Sleep(CrawlerInterval)
		log.Debug("tv id: ", id)
		tv, err := tu.providerRepo.GetTVDetail(id)
		if err != nil {
			if _, ok := err.(provider.APINotFoundError); !ok {
				log.WithError(err).Error("Get discover Fail")
			}
			continue
		}
		for i, _ := range tv.Seasons {
			// time.Sleep(CrawlerInterval)
			avg, count, err := tu.providerRepo.GetTVSeasonVote(uint(id), tv.Seasons[i].SeasonNumber)
			if err != nil {
				log.WithError(err).Error("get tv season vote fail")
				continue
			}
			tv.Seasons[i].VoteAverage = avg
			tv.Seasons[i].VoteCount = count
		}
		ch <- tv
	}
	return
}

var tvWriteChan chan *models.TV
var tvWriteSyncOnce sync.Once

func (tu *tmdbUsecase) CreateStoreTVTask() chan *models.TV {
	tvWriteSyncOnce.Do(func() {
		tvWriteChan = make(chan *models.TV, 100000)
		go func() {
			for {
				tv := <-tvWriteChan
				if _, err := tu.movieRepo.TVStore(tv); err != nil {
					log.WithField("tv", tv).WithError(err).Error("tv Task store fail")
				}
			}
		}()
	})

	return tvWriteChan
}

func (tu *tmdbUsecase) StartCrawlerCredit(ch chan *models.Credit, pch chan *models.Person) {
	order := map[string]string{"popularity": "desc"}
	p, totalPage := 1, 1
	for {
		movies, page, err := tu.movieRepo.MovieList(1, 100, order, nil)
		if err != nil {
			log.WithError(err).Error("Crealer credit fail on get movies id")
		}
		if p == 1 {
			totalPage = int(page.TotalPages)
		}
		for _, m := range movies {
			casts, crews, err := tu.providerRepo.GetMovieCredits(m.ProviderID)
			if err != nil {
				log.WithError(err).WithField("movie", m).Error("get credits fail")
			}
			for i := range casts {
				if casts[i].Order != nil && *casts[i].Order > 5 {
					continue
				}
				providerPersonID := casts[i].ProviderPersonID
				id, err := tu.movieRepo.PeopleIDByProviderID(providerPersonID)
				if err != nil {
					logfield := log.Fields{
						"provider People ID": providerPersonID,
						"provider movie ID":  m.ProviderID,
					}
					if err != gorm.ErrRecordNotFound {
						log.WithFields(logfield).WithError(err).Error("find people id from db fail")
						continue
					}
					p, err := tu.providerRepo.GetPersonDetail(int(providerPersonID))
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("get people from provider fail")
						continue
					}
					id, err = tu.personRepo.Store(p)
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("store people fail")
						continue
					}
				}
				casts[i].PersonID = id
				ch <- &casts[i]
			}
			for i := range crews {
				providerPersonID := crews[i].ProviderPersonID
				id, err := tu.movieRepo.PeopleIDByProviderID(providerPersonID)
				if err != nil {
					logfield := log.Fields{
						"provider People ID": providerPersonID,
						"provider movie ID":  m.ProviderID,
					}
					if err != gorm.ErrRecordNotFound {
						log.WithFields(logfield).WithError(err).Error("find people id from db fail")
						continue
					}
					p, err := tu.providerRepo.GetPersonDetail(int(providerPersonID))
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("get people from provider fail")
						continue
					}
					id, err = tu.personRepo.Store(p)
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("store people fail")
						continue
					}
				}
				crews[i].PersonID = id
				ch <- &crews[i]
			}
		}

		if p > totalPage {
			break
		}
		p++
	}
	return
}

func (tu *tmdbUsecase) StartCrawlerPopularMovie(ch chan *models.Credit, pch chan *models.Person, mch chan *models.Movie) {
	total, err := tu.providerRepo.GetMovieTotalPages()
	if err != nil {
		log.WithError(err).Error("get total page fail")
		return
	}
	opt := map[string]string{}
	opt["sort_by"] = "popularity.desc"
	log.Info("start crawler Popular Movie")
	for i := 1; i <= total; i++ {
		if i > 1000 {
			break
		}
		movies, err := tu.providerRepo.GetMovieWithPage(i, opt)
		if err != nil {
			log.WithField("page", i).WithError(err).Error("get movies by page fail")
			continue
		}
		for _, m := range movies {
			mch <- m
			casts, crews, err := tu.providerRepo.GetMovieCredits(m.ProviderID)
			if err != nil {
				log.WithError(err).WithField("movie", m).Error("get credits fail")
			}
			for i := range casts {
				if casts[i].Order != nil && *casts[i].Order > 5 {
					continue
				}
				providerPersonID := casts[i].ProviderPersonID
				id, err := tu.movieRepo.PeopleIDByProviderID(providerPersonID)
				if err != nil {
					logfield := log.Fields{
						"provider People ID": providerPersonID,
						"provider movie ID":  m.ProviderID,
					}
					if err != gorm.ErrRecordNotFound {
						log.WithFields(logfield).WithError(err).Error("find people id from db fail")
						continue
					}
					p, err := tu.providerRepo.GetPersonDetail(int(providerPersonID))
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("get people from provider fail")
						continue
					}
					id, err = tu.personRepo.Store(p)
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("store people fail")
						continue
					}
				}
				casts[i].PersonID = id
				ch <- &casts[i]
			}
			for i := range crews {
				providerPersonID := crews[i].ProviderPersonID
				id, err := tu.movieRepo.PeopleIDByProviderID(providerPersonID)
				if err != nil {
					logfield := log.Fields{
						"provider People ID": providerPersonID,
						"provider movie ID":  m.ProviderID,
					}
					if err != gorm.ErrRecordNotFound {
						log.WithFields(logfield).WithError(err).Error("find people id from db fail")
						continue
					}
					p, err := tu.providerRepo.GetPersonDetail(int(providerPersonID))
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("get people from provider fail")
						continue
					}
					id, err = tu.personRepo.Store(p)
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("store people fail")
						continue
					}
				}
				crews[i].PersonID = id
				ch <- &crews[i]
			}
		}
	}
	log.Info("end crawler Popular Movie")
}

func (tu *tmdbUsecase) StartCrawlerPopularTV(ch chan *models.Credit, pch chan *models.Person, tvch chan *models.TV) {
	total, err := tu.providerRepo.GetTVTotalPages()
	if err != nil {
		log.WithError(err).Error("get total page fail")
		return
	}
	opt := map[string]string{}
	opt["sort_by"] = "popularity.desc"
	log.Info("start crawler Popular TV")
	for i := 1; i <= total; i++ {
		if i > 1000 {
			break
		}
		tvData, err := tu.providerRepo.GetTVWithPage(i, opt)
		if err != nil {
			log.WithField("page", i).WithError(err).Error("get tvData by page fail")
			continue
		}
		for _, tv := range tvData {
			tvDetail, err := tu.providerRepo.GetTVDetail(int(tv.ProviderID))
			if err != nil {
				if _, ok := err.(provider.APINotFoundError); !ok {
					log.WithError(err).Error("Get discover Fail")
				}
				continue
			}
			for i, _ := range tvDetail.Seasons {
				avg, count, err := tu.providerRepo.GetTVSeasonVote(tv.ProviderID, tvDetail.Seasons[i].SeasonNumber)
				if err != nil {
					log.WithError(err).Error("get tv season vote fail")
					continue
				}
				tvDetail.Seasons[i].VoteAverage = avg
				tvDetail.Seasons[i].VoteCount = count
			}
			tvch <- tvDetail
			casts, crews, err := tu.providerRepo.GetTVCredits(tvDetail.ProviderID)
			if err != nil {
				log.WithError(err).WithField("tv", tvDetail).Error("get tv credits fail")
			}
			for i := range casts {
				if casts[i].Order != nil && *casts[i].Order > 5 {
					continue
				}
				providerPersonID := casts[i].ProviderPersonID
				id, err := tu.movieRepo.PeopleIDByProviderID(providerPersonID)
				if err != nil {
					logfield := log.Fields{
						"provider People ID": providerPersonID,
						"provider tv ID":     tvDetail.ProviderID,
					}
					if err != gorm.ErrRecordNotFound {
						log.WithFields(logfield).WithError(err).Error("find people id from db fail")
						continue
					}
					p, err := tu.providerRepo.GetPersonDetail(int(providerPersonID))
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("get people from provider fail")
						continue
					}
					id, err = tu.personRepo.Store(p)
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("store people fail")
						continue
					}
				}
				casts[i].PersonID = id
				ch <- &casts[i]
			}
			for i := range crews {
				providerPersonID := crews[i].ProviderPersonID
				id, err := tu.movieRepo.PeopleIDByProviderID(providerPersonID)
				if err != nil {
					logfield := log.Fields{
						"provider People ID": providerPersonID,
						"provider tv ID":     tvDetail.ProviderID,
					}
					if err != gorm.ErrRecordNotFound {
						log.WithFields(logfield).WithError(err).Error("find people id from db fail")
						continue
					}
					p, err := tu.providerRepo.GetPersonDetail(int(providerPersonID))
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("get people from provider fail")
						continue
					}
					id, err = tu.personRepo.Store(p)
					if err != nil {
						log.WithError(err).WithFields(logfield).Error("store people fail")
						continue
					}
				}
				crews[i].PersonID = id
				ch <- &crews[i]
			}
		}
	}
	log.Info("end crawler Popular TV")
}

var creditWriteChan chan *models.Credit
var creditWriteSyncOnce sync.Once

func (tu *tmdbUsecase) CreateStoreCreditTask() chan *models.Credit {
	creditWriteSyncOnce.Do(func() {
		creditWriteChan = make(chan *models.Credit, 100000)
		go func() {
			for {
				c := <-creditWriteChan
				if _, err := tu.movieRepo.CreditStore(c); err != nil {
					log.WithError(err).Error("credit Task store fail")
				}
			}
		}()
	})

	return creditWriteChan
}
