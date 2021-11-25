package main

import (
	"fmt"
	twitterstream "github.com/fallenstedt/twitter-stream"
	"time"
)

const key = "YOUR_KEY"
const secret = "YOUR_SECRET"

func main() {

 	addRules()
 	getRules()
 	// You can delete the rules created in this example
 	//deleteRules()
}

type StreamData struct {
	Data struct {
		Text      string    `json:"text"`
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		AuthorID  string    `json:"author_id"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
	MatchingRules []struct {
		ID  string  `json:"id"`
		Tag string `json:"tag"`
	} `json:"matching_rules"`
}


func addRules() {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(key, secret).RequestBearerToken()
	if err != nil {
		panic(err)
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken)
	res, err := api.Rules.AddRules(`{
		"add": [
				{"value": "cat has:images", "tag": "cat tweets with images"}
			]
		}`, false) // dryRun is set to false.

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twitter: %v", res.Errors))
	}

	fmt.Println("I have created this many rules: ")
	fmt.Println(res.Meta.Summary.Created)
}

func getRules() {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(key, secret).RequestBearerToken()
	if err != nil {
		panic(err)
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken)
	res, err := api.Rules.GetRules()

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twitter: %v", res.Errors))
	}

	fmt.Println(res.Data)
}


func deleteRules() {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(key, secret).RequestBearerToken()
	if err != nil {
		panic(err)
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken)

	// use api.Rules.GetRules to find the ID number for an existing rule
	res, err := api.Rules.AddRules(`{
		"delete": {
				"ids": ["1234567890"]
			}
		}`, false)

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twitter: %v", res.Errors))
	}

	fmt.Println(res)
}