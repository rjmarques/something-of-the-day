package service
import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

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
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/something", sm.getSomethingHandler)

	// TODO PUBLISH UI

	compressedRouter := handlers.CompressHandler(mux)
	loggedRouter := handlers.LoggingHandler(os.Stdout, compressedRouter)
	log.Fatal(http.ListenAndServe(":80", loggedRouter))
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
			var sinceId int64 = firstRelevantTweetID
			if latestSomething := sm.st.GetLatest(); latestSomething != nil {
				log.Printf("didn't found a something in the store, defaulting to something ID %d", firstRelevantTweetID)
				sinceId = latestSomething.Id
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

func (sm *SomethingOfTheDay) getSomethingHandler(w http.ResponseWriter, r *http.Request) {
	something := sm.st.GetRand()
	json.NewEncoder(w).Encode(something)
}

