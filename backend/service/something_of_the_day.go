package service

import (
	"log"
	"net/http"
	"time"

	"github.com/rjmarques/something-of-the-day/datastore"
	"github.com/rjmarques/something-of-the-day/model"
	"github.com/rjmarques/something-of-the-day/scraper"
	"github.com/rjmarques/something-of-the-day/store"
)

const (
	firstRelevantTweetID = 849296076343582720
)

type SomethingOfTheDay struct {
	dao *datastore.DAO
	st  *store.Store
	sc  *scraper.Scraper
}

func NewSomethingOfTheDay(clientID, clientSecret string, dao *datastore.DAO) *SomethingOfTheDay {
	return &SomethingOfTheDay{
		dao: dao,
		sc:  scraper.NewScraper(clientID, clientSecret),
		st:  store.NewStore(),
	}
}

func (sm *SomethingOfTheDay) GetRandomSomething() *model.Something {
	return sm.st.GetRand()
}

func (sm *SomethingOfTheDay) Start() {
	err := sm.hydrateStore()
	if err != nil {
		panic(err)
	}

	// background routine that picks up new somethings of the day
	go sm.update()

	// Rest API endpoints
	api := apiServer(sm.st)
	log.Fatal(http.ListenAndServe(":80", api))
}

func (sm *SomethingOfTheDay) hydrateStore() error {
	log.Println("hydrating store")

	somethings, err := sm.dao.GetSomethings()
	if err != nil {
		return err
	}

	log.Printf("found %d somethings in the datastore", len(somethings))
	sm.st.Add(somethings)

	return nil
}

func (sm *SomethingOfTheDay) update() {
	t := time.Tick(1 * time.Minute)
	for {
		select {
		case <-t:
			var sinceId int64
			if latestSomething := sm.st.GetLatest(); latestSomething != nil {
				sinceId = latestSomething.Id
			} else {
				log.Printf("didn't find a something in the store, defaulting to something ID %d", firstRelevantTweetID)
				sinceId = firstRelevantTweetID
			}

			somethings, err := sm.sc.GetNewerSomethings(sinceId)
			if err != nil {
				log.Printf("failed to get somethings: %v", err)
			}

			log.Printf("found %d new somethings", len(somethings))

			err = sm.dao.StoreSomethings(somethings)
			if err != nil {
				log.Printf("failed to store new somethings: %v", err)
			} else {
				sm.st.Add(somethings)
			}
		}
	}
}
