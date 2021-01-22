package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"database/sql"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Raja@1234"
	dbname   = "tweetbot"
)

type Api struct {
	consumerkey    string
	consumersecret string
	accesstoken    string
	accesssecret   string
}

func dbcheck() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB Connection Successful")
}

func connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	results, err := db.Query("SELECT consumerkey, consumersecret, accesstoken, accesssecret FROM apitokens")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var api Api
		err = results.Scan(&api.consumerkey, &api.consumersecret, &api.accesstoken, &api.accesssecret)
		if err != nil {
			panic(err.Error())
		}
		consumerKey := (api.consumerkey)
		consumerSecret := (api.consumersecret)
		accessToken := (api.accesstoken)
		accessSecret := (api.accesssecret)

		config := oauth1.NewConfig(consumerKey, consumerSecret)
		token := oauth1.NewToken(accessToken, accessSecret)

		//OAuth1 http.Client will automatically authorize Requests
		httpClient := config.Client(oauth1.NoContext, token)

		// Twitter client
		client := twitter.NewClient(httpClient)

		// Verify Credentials
		verifyParams := &twitter.AccountVerifyParams{
			IncludeEmail: twitter.Bool(true),
		}
		user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
		fmt.Printf("\nTwitter Bot Initialized...\n\n")
		currentTime := time.Now()
		fmt.Printf("User's Name: %+v\t %v\n\n", user.ScreenName, currentTime.Format("01-02-2006 15:04:05"))
		fmt.Printf("Hello, %v\n\n", user.Name)

		fmt.Printf("Options:\n\n")
		fmt.Println("1. Home Timeline")
		fmt.Println("2. Tweet")
		fmt.Println("3. Status Check")
		fmt.Println("4. Search Tweets")
		fmt.Println("5. Show User")
		fmt.Printf("6. Followers\n\n")
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Enter your option: ")
		scanner.Scan()
		option, _ := strconv.ParseInt(scanner.Text(), 10, 64)

		switch option {
		case 1:
			tweets, resp, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
				Count: 20,
			})
			if err != nil {
				log.Println(err)
			}
			log.Printf("%v\n", resp.Body)
			log.Printf("%v\n", tweets)
			break

		case 2:
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("Enter your tweet: ")
			scanner.Scan()
			tweetinput := scanner.Text()

			fmt.Printf("Your Tweet: %s", tweetinput)

			tweet, resp, err := client.Statuses.Update(tweetinput, nil)
			if err != nil {
				log.Println(err)
			}
			log.Printf("\n\n%+v\n", resp)
			log.Printf("%+v\n", tweet)
			break

		case 3:
			tweet, resp, err := client.Statuses.Show(585613041028431872, nil)
			if err != nil {
				log.Println(err)
			}
			log.Printf("\n\n%+v\n", resp)
			log.Printf("%+v\n", tweet)
			break

		case 4:
			search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
				Query: "Hello",
			})
			if err != nil {
				log.Println(err)
			}

			log.Printf("\n\n%v\n", resp.Status)
			log.Printf("%v\n", search.Statuses)
			break

		case 5:
			user, resp, err := client.Users.Show(&twitter.UserShowParams{
				ScreenName: "dghubble"})
			if err != nil {
				log.Println(err)
			}
			log.Printf("\n\n%+v\n", resp)
			log.Printf("%+v\n", user)
			break

		case 6:
			followers, resp, err := client.Followers.List(&twitter.FollowerListParams{})
			if err != nil {
				log.Println(err)
			}
			log.Printf("\n\n%+v\n", resp.Body)
			log.Printf("%+v\n", followers)
			break

		default:
			fmt.Println("Incorrect Option")
			break
		}
	}
}

func main() {

	dbcheck()

	connect()

}
