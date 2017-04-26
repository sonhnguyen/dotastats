package dotastats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/mrjones/oauth"
	"github.com/spf13/viper"
)

const (
	CreateListURL     = "https://api.twitter.com/1.1/lists/create.json"
	AddMembersURL     = "https://api.twitter.com/1.1/lists/members/create_all.json"
	RemoveListURL     = "https://api.twitter.com/1.1/lists/destroy.json"
	RequestTokenUrl   = "https://api.twitter.com/oauth/request_token"
	AuthorizeTokenUrl = "https://api.twitter.com/oauth/authorize"
	AccessTokenUrl    = "https://api.twitter.com/oauth/access_token"
)

func CreateOAuth() (*http.Client, error) {
	var consumerKey, consumerSecret, accessToken, accessTokenSecret string

	if viper.GetBool("isDevelopment") {
		twitterCred := viper.GetStringMapString("twitter")
		consumerKey = twitterCred["consumerKey"]
		consumerSecret = twitterCred["consumerSecret"]
		accessToken = twitterCred["token"]
		accessTokenSecret = twitterCred["tokenSecret"]
	} else {
		consumerKey = os.Getenv("consumerKey")
		consumerSecret = os.Getenv("consumerSecret")
		accessToken = os.Getenv("token")
		accessTokenSecret = os.Getenv("tokenSecret")
	}

	c := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   RequestTokenUrl,
			AuthorizeTokenUrl: AuthorizeTokenUrl,
			AccessTokenUrl:    AccessTokenUrl,
		})

	t := oauth.AccessToken{
		Token:  accessToken,
		Secret: accessTokenSecret,
	}
	return c.MakeHttpClient(&t)
}

func AddMembersToListTwitter(client *http.Client, req TwitterAddToListRequest) error {
	response, err := client.PostForm(AddMembersURL,
		url.Values{
			"owner_screen_name": []string{req.OwnerScreenName},
			"slug":              []string{req.Slug},
			"screen_name":       []string{req.ScreenName},
		})

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("error on adding member to twitter list, %s, %s, %s\n", req.Slug, req.ScreenName, body)
	}
	if err != nil {
		return err
	}
	return nil
}

func CreateListTwitter(client *http.Client, req TwitterCreateListRequest) error {
	response, err := client.PostForm(CreateListURL,
		url.Values{
			"name":        []string{req.Name},
			"mode":        []string{req.Mode},
			"description": []string{req.Description},
		})

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("error on creating twitter list, %s, %s\n", req.Name, body)
	}
	if req.Name == "dota-team-np" {
		fmt.Printf("%s\n", body)
	}
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllListFromTwitter(client *http.Client, twitterID string) error {
	response, err := client.Get("https://api.twitter.com/1.1/lists/ownerships.json?screen_name=dotastats_&count=800")

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("error on getting all twitter list\n")
	}
	twitterGetListResponse := TwitterGetListResponse{}

	_ = json.Unmarshal(body, &twitterGetListResponse)

	for _, list := range twitterGetListResponse.Lists {
		_ = RemoveListFromTwitter(client, TwitterRemoveListRequest{
			OwnerScreenName: "dotastats_",
			Slug:            list.Slug,
		})
	}

	return nil
}

func RemoveListFromTwitter(client *http.Client, req TwitterRemoveListRequest) error {
	response, err := client.PostForm(RemoveListURL,
		url.Values{
			"owner_screen_name": []string{req.OwnerScreenName},
			"slug":              []string{req.Slug},
		})

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("error on removing twitter list, %s, %s\n", req.Slug, body)
	}
	if err != nil {
		return err
	}
	return nil
}
