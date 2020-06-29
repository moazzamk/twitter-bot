package twitter_bot

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {
	config := &clientcredentials.Config{
		ClientID: "consumerKey",
		ClientSecret: "secret",
		TokenURL: "https://api.twitter.com/oauth2/token",
	}


	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	client
	fmt.Println(client)
}