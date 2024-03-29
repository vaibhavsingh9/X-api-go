package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/oauth1"
)

// Struct to parse webhook load
type WebhookLoad struct {
	UserId           string  `json:"for_user_id"`
	TweetCreateEvent []Tweet `json:"tweet_create_events"`
}

// Struct to parse tweet
type Tweet struct {
	Id    int64
	IdStr string `json:"id_str"`
	User  User
	Text  string
}

// Struct to parse user
type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
}

func CreateClient() *http.Client {
	//Create oauth client with consumer keys and access token
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))

	subscribeWebhook()
	return config.Client(oauth1.NoContext, token)

}

func registerWebhook() {
	fmt.Println("Registering webhook...")
	httpClient := CreateClient()

	//Set parameters
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/webhooks.json"
	values := url.Values{}
	values.Set("url", os.Getenv("APP_URL")+"/webhook/twitter")

	//Make Oauth Post with parameters
	resp, _ := httpClient.PostForm(path, values)
	defer resp.Body.Close()
	//Parse response and check response
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		panic(err)
	}
	fmt.Println("Webhook id of " + data["id"].(string) + " has been registered")
}

func subscribeWebhook() {
	fmt.Println("subscribing webapp...")
	client := CreateClient()
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/subscriptions.json"
	resp, _ := client.PostForm(path, nil)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		fmt.Println("subscribed successfully")
	} else if resp.StatusCode != 204 {
		fmt.Println(("could not subscribe the webhook"))
		fmt.Println(string(body))
	}
}
