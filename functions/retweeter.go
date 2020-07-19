package functions

import (
	"context"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"math"
	"math/rand"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func Retweeter(ctx context.Context, m PubSubMessage) error {
	config := oauth1.NewConfig("d", "d")
	token := oauth1.NewToken("d", "d")

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	tags := []string{
		"#hyperrealistic",
		"#realism",
		"#monet",
		"#russianart #realism",
		"#americanart #realism",
		"#ivanshishkin",
		"@Albert Bierstadt",
	}

	length := len(tags)
	key := rand.Intn(length - 1)

	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query:  tags[key],
		Count:  5,
		Lang:   "en",
		Filter: "images",
	})

	if err != nil {
		log.Println("ERROR: Search failed", err)
	}

	length = int(math.Min(
		float64(len(search.Statuses)),
		10,
	))
	tweetKey := rand.Intn(length)

	log.Println("KEY: ", key, "VALUE: ", tags[key], "TWEET_KEY: ", tweetKey)
	log.Println(search.Statuses[tweetKey].Text)

	_, _, err = client.Statuses.Retweet(
		search.Statuses[tweetKey].ID,
		&twitter.StatusRetweetParams{},
	)

	if err != nil {
		log.Println("ERROR: Retweet failed", err)
	}

	return nil
}
