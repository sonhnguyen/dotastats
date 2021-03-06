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
	RateLimitURL      = "https://api.twitter.com/1.1/application/rate_limit_status.json"
	RemoveListURL     = "https://api.twitter.com/1.1/lists/destroy.json"
	RequestTokenUrl   = "https://api.twitter.com/oauth/request_token"
	AuthorizeTokenUrl = "https://api.twitter.com/oauth/authorize"
	AccessTokenUrl    = "https://api.twitter.com/oauth/access_token"
)

func CreateOAuth() (*http.Client, error) {
	var consumerKey, consumerSecret, accessToken, accessTokenSecret string

	if viper.GetBool("isDevelopment") {
		twitterCred := viper.GetStringMapString("twitter")
		consumerKey = twitterCred["consumerkey"]
		consumerSecret = twitterCred["consumersecret"]
		accessToken = twitterCred["token"]
		accessTokenSecret = twitterCred["tokensecret"]
	} else {
		consumerKey = os.Getenv("consumerKey")
		consumerSecret = os.Getenv("consumersecret")
		accessToken = os.Getenv("token")
		accessTokenSecret = os.Getenv("tokensecret")
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

func CheckTwitterRateLimit(client *http.Client) error {
	response, err := client.Get("https://api.twitter.com/1.1/application/rate_limit_status.json?resources=lists")

	if err != nil {
		fmt.Printf("error on post form, %s\n", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("error on checking twitter rate limit, %s\n", body)
	}
	fmt.Println(string(body))
	if err != nil {
		return err
	}
	return nil
}

func AddMembersToListTwitter(client *http.Client, req TwitterAddToListRequest) error {
	response, err := client.PostForm(AddMembersURL,
		url.Values{
			"owner_screen_name": []string{req.OwnerScreenName},
			"slug":              []string{req.Slug},
			"screen_name":       []string{req.ScreenName},
		})

	if err != nil {
		fmt.Printf("error on post form, %s\n", err)
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
		fmt.Printf("error on post form, %s\n", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("error on creating twitter list, %s, %s\n", req.Name, body)
	}
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllListFromTwitter(twitterID string) error {
	client, err := CreateOAuth()
	if err != nil {
		return err
	}
	response, err := client.Get("https://api.twitter.com/1.1/lists/ownerships.json?screen_name=dotastats_&count=800")

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("error on getting all twitter list, %s\n,", body)
	}
	twitterGetListResponse := TwitterGetListResponse{}

	_ = json.Unmarshal(body, &twitterGetListResponse)

	for _, list := range twitterGetListResponse.Lists {
		_ = RemoveListFromTwitter(client, TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
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
		fmt.Printf("error on post form, %s\n", err)
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

func CreateTwitterList(teams []TeamInfo) error {
	var errorList []error
	c, err := CreateOAuth()
	if err != nil {
		return err
	}
	var twitterID string
	if viper.GetBool("isDevelopment") {
		twitterID = viper.GetString("twitter.twitterID")
	} else {
		twitterID = os.Getenv("twitterid")
	}

	for _, team := range teams {
		nameSlug := team.Game + "-" + team.NameSlug
		if len(nameSlug) > 25 {
			nameSlug = nameSlug[:25]
		}

		err = RemoveListFromTwitter(c, TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug,
		})

		if err != nil {
			errorList = append(errorList, err)
		}

		err = CreateListTwitter(c, TwitterCreateListRequest{
			Name:        nameSlug,
			Mode:        "public",
			Description: team.Game + " - " + team.Region + " - " + team.Name,
		})

		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		memberScreenNames := ""
		for _, player := range team.Players {
			screenName := player.FindTwitterID()
			if len(screenName) == 0 {
				continue
			}
			memberScreenNames += screenName + ","
		}
		if memberScreenNames == "" {
			continue
		}
		memberScreenNames = memberScreenNames[:len(memberScreenNames)-1]
		err = AddMembersToListTwitter(c, TwitterAddToListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug,
			ScreenName:      memberScreenNames,
		})

		if err != nil {
			errorList = append(errorList, err)
		}
	}

	return fmt.Errorf("error when save team list to twitter", errorList)
}
