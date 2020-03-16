package scraper
import (
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/rjmarques/something-of-the-day/model"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Scraper struct {
	client *twitter.Client
}

func NewScraper(clientID, clientSecret string) *Scraper {
	return &Scraper{
		client: initClient(clientID, clientSecret),
	}
}

func initClient(clientID, clientSecret string) *twitter.Client {
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return client
}

// Returns Something's in order from oldest to newest
func (sc *Scraper) GetNewerSomethings(sinceID int64) ([]*model.Something, error) {
	log.Printf("getting new somethings, since ID %d", sinceID)

	excludeReplies := true
	includeRetweets := false
	params := &twitter.UserTimelineParams{
		ScreenName:      "lgst_something",
		Count:           50,
		IncludeRetweets: &includeRetweets,
		ExcludeReplies:  &excludeReplies,
		TweetMode:       "extended",
		SinceID:         sinceID,
	}

	somethings := []*model.Something{}
	for {
		// user show
		tweets, resp, err := sc.client.Timelines.UserTimeline(params)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("error in response: %v", resp)
		}

		if len(tweets) == 0 {
			break
		}

		for _, t := range tweets {
			createdAt, err := time.Parse(time.RubyDate, t.CreatedAt)
			if err != nil {
				return nil, err
			}
			somethings = append(somethings, &model.Something{
				Id:        t.ID,
				CreatedAt: createdAt,
				Text:      t.FullText,
			})
		}

		maxId := tweets[len(tweets)-1].ID
		if params.MaxID == maxId {
			break
		}
		params.MaxID = maxId - 1 // as suggested in https://developer.twitter.com/en/docs/tweets/timelines/guides/working-with-timelines
	}

	return inverse(somethings), nil
}

func inverse(list []*model.Something) []*model.Something {
	if len(list) <= 1 {
		return list
	}

	var inversed []*model.Something
	for i := len(list) - 1; i >= 0; i-- {
		inversed = append(inversed, list[i])
	}
	return inversed
}

