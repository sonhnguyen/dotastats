package dotastats

import (
	"fmt"
	"json"

	"github.com/mrjones/oauth"
	"github.com/spf13/viper"
)

const (
	CreateListURL     = "https://api.twitter.com/1.1/lists/create.json"
	AddMemberURL      = "https://api.twitter.com/1.1/lists/members/create.json"
	RemoveMemberURL   = "https://api.twitter.com/1.1/lists/members/destroy_all.json"
	GetMemberURL      = "https://api.twitter.com/1.1/lists/members.json"
	RequestTokenUrl   = "https://api.twitter.com/oauth/request_token"
	AuthorizeTokenUrl = "https://api.twitter.com/oauth/authorize"
	AccessTokenUrl    = "https://api.twitter.com/oauth/access_token"
)

func CreateOAuth() (*http.Client, error) {
	if viper.GetBool("isDevelopment") {
		twitterCred := viper.GetStringMapString("twitter")
		consumerKey := twitterCred["ConsumerKey"]
		consumerSecret := twitterCred["ConsumerSecret"]
		accessToken := twitterCred["AccessToken"]
		accessTokenSecret := twitterCred["AccessTokenSecret"]
	} else {
		consumerKey := os.Getenv("ConsumerKey")
		consumerSecret := os.Getenv("ConsumerSecret")
		accessToken := os.Getenv("AccessToken")
		accessTokenSecret := os.Getenv("AccessTokenSecret")
	}

	if !consumerKey || !consumerSecret || !accessToken || !accessTokenSecret {
		return fmt.Errorf("error on getting twitter credentials")
	}

	c := oauth.NewConsumer(
		*consumerKey,
		*consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   RequestTokenUrl,
			AuthorizeTokenUrl: AuthorizeTokenUrl,
			AccessTokenUrl:    AccessTokenUrl,
		})

	t := oauth.AccessToken{
		Token:  *accessToken,
		Secret: *accessTokenSecret,
	}
	return c.MakeHttpClient(&t)
}

func AddMemberToListTwitter(client *http.Client, req TwitterAddToListRequest) error {
	response, err := client.Post(AddMemberURL,
		url.Values{
			"owner_screen_name": []string{req.OwnerScreenName},
			"slug":              []string{req.Slug},
			"screen_name":       []string{req.ScreenName},
		})

	if err != nil {
		return err
	}

	bits, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("error on creating twitter list, %s", bits)
	}

	return nil
}

func DeleteMemberFromListTwitter(client *http.Client, req TwitterRemoveFromListRequest) error {
	response, err := client.Post(RemoveMemberURL,
		url.Values{
			"owner_screen_name": []string{req.OwnerScreenName},
			"slug":              []string{req.Slug},
			"screen_name":       []string{req.ScreenName},
		})

	if err != nil {
		return err
	}

	bits, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("error on creating twitter list, %s", bits)
	}

	return nil
}

func CreateListTwitter(client *http.Client, req TwitterCreateListRequest) error {
	response, err := client.Post(CreateListURL,
		url.Values{
			"name":        []string{req.Name},
			"mode":        []string{req.Mode},
			"description": []string{req.Description},
		})

	if err != nil {
		return err
	}

	bits, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("error on creating twitter list, %s", bits)
	}

	return nil
}

func GetMemberFromListTwitter(client *http.Client, req TwitterGetFromListRequest) ([]TwitterUser, error) {
	response, err := client.Post(GetMemberURL,
		url.Values{
			"owner_screen_name": []string{req.OwnerScreenName},
			"slug":              []string{req.Slug},
		})

	if err != nil {
		return err
	}

	bits, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("error on creating twitter list, %s", bits)
	}
	twitterGetFromListResponse := TwitterGetFromListResponse{}
	err = json.Unmarshal(bits, &twitterGetFromListResponse)

	if err != nil {
		return twitterGetFromListResponse, err
	}
	return twitterGetFromListResponse, nil
}
