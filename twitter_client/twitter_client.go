package twitter_client

import (
	"../configuration"
	"context"
	"fmt"
	_ "github.com/creasty/defaults"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type TwitterClient struct {
	Config 			*configuration.Configuration	`default:"-"`
	httpClient 		*http.Client					`default:"-"`
	client 			*twitter.Client					`default:"-"`
	DefaultCount 	int 							`default:"20"`
}

func (t *TwitterClient) Login() {
	if t.Config.Key == "" {
		fmt.Print("Please provide key and token in ~/.twitter_cli")
		os.Exit(-1)
	}

	var config = &oauth1.Config {
		ConsumerKey:    t.Config.Key,
		ConsumerSecret: t.Config.Token,
		CallbackURL:    "http://localhost:6723/callback",
		Endpoint:       oauth1.Endpoint{RequestTokenURL:"https://api.twitter.com/oauth/request_token", AuthorizeURL:"https://api.twitter.com/oauth/authorize", AccessTokenURL:"https://api.twitter.com/oauth/access_token"},
	}

	if t.Config.AccessToken == "" {
		var requestSecret = ""
		logger := log.New(os.Stdout, "http: ", log.LstdFlags)
		logger.Println("Server is starting...")

		done := make(chan bool)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		router := http.NewServeMux()

		server := &http.Server{
			Addr:         ":6723",
			Handler:      router,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}

		router.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestToken = "";
			requestToken, requestSecret, _ = config.RequestToken()
			authorizationURL, _ := config.AuthorizationURL(requestToken)
			// handle err
			http.Redirect(w, r, authorizationURL.String(), http.StatusFound)
		}))

		router.Handle("/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestToken, verifier, _ := oauth1.ParseAuthorizationCallback(r)

			accessToken, accessSecret, _ := config.AccessToken(requestToken, requestSecret, verifier)
			// handle error
			t.Config.AccessToken = accessToken
			t.Config.AccessSecret = accessSecret
			t.Config.SaveConfiguration()

			close(quit)
		}))

		go func() {
			<-quit
			logger.Println("Server is shutting down...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			server.SetKeepAlivesEnabled(false)
			if err := server.Shutdown(ctx); err != nil {
				logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
			}
			close(done)
		}()
		fmt.Print("Please open http://localhost:6723/login in your browser")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on %s: %v\n", ":6723", err)
		}

		<-done
	}

	token := oauth1.NewToken(t.Config.AccessToken, t.Config.AccessSecret)
	t.httpClient = config.Client(oauth1.NoContext, token)
	t.client = twitter.NewClient(t.httpClient)
}

func (t *TwitterClient) GetTweets(count int32) []twitter.Tweet {
	params := &twitter.HomeTimelineParams{
		Count: t.DefaultCount,
		TweetMode:"extended",
	}
	tweets, _, _ := t.client.Timelines.HomeTimeline(params)
	return tweets
}

func (t *TwitterClient) GetTweetsFromUser(user string) []twitter.Tweet {
	params := &twitter.UserTimelineParams{
		Count: t.DefaultCount,
		TweetMode:"extended",
		ScreenName:user,
	}
	tweets, _, _ := t.client.Timelines.UserTimeline(params)
	return tweets
}

func (t *TwitterClient) GetMentions() []twitter.Tweet {
	params := &twitter.MentionTimelineParams{
		Count: t.DefaultCount,
		TweetMode:"extended",
	}
	tweets, _, _ := t.client.Timelines.MentionTimeline(params)
	return tweets
}

func (t *TwitterClient) GetRateLimits() *twitter.RateLimit {
	limit, _, _ := t.client.RateLimits.Status(&twitter.RateLimitParams{})
	return limit
}

func (t *TwitterClient) GetTrendLocations() []twitter.Location {
	locations, _, _ := t.client.Trends.Available()
	return locations
}

func (t *TwitterClient) GetFollowers(user string) []twitter.User {
	params := &twitter.FollowerListParams{
		Count: t.DefaultCount,
		ScreenName: user,
	}
	var result []twitter.User
	for {
		followers, _, _ := t.client.Followers.List(params)
		params.Cursor = followers.NextCursor
		result = append(result, followers.Users...)
		if len(followers.NextCursorStr) <= 0 || len(result) >= t.DefaultCount {
			break;
		}
	}
	return result
}