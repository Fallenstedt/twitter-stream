package main

import (
	"encoding/json"
	"fmt"
	twitterstream "github.com/fallenstedt/twitter-stream"
	"log"
	"time"
)

const key = "YOUR_KEY"
const secret = "YOUR_SECRET"


func main() {
	// Use your favorite function from below here
 	startStream()
 	addRules()
 	getRules()
	deleteRules()
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
		ID  int64  `json:"id"`
		Tag string `json:"tag"`
	} `json:"matching_rules"`
}



func startStream() {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(key, secret).RequestBearerToken()

	if err != nil {
		panic(err)
	}

	api := twitterstream.NewTwitterStream(tok.AccessToken)
	// It is encouraged you unmarashal json with twitterstream's unmarshal hook. This is a thread-safe
	// way to unmarshal json
	api.Stream.SetUnmarshalHook(func(bytes []byte) (interface{}, error) {
		data := StreamData{}
		if err := json.Unmarshal(bytes, &data); err != nil {
			log.Printf("Failed to unmarshal bytes: %v", err)
		}
		return data, err
	})
	defer  api.Stream.StopStream()

	err = api.Stream.StartStream("?expansions=author_id&tweet.fields=created_at")

	if err != nil {
		panic(err)
	}

	go func() {
		for message := range api.Stream.GetMessages() {
			if message.Err != nil {
				panic(message.Err)
			}

			// Type assertion
			tweet, ok := message.Data.(StreamData)
			if !ok {
				continue
			}

			fmt.Println(tweet.Data.Text)
		}
	}()

	time.Sleep(time.Second * 30)
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
		}`, true) // dryRun is set to true

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twiiter: %v", res.Errors))
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
		panic(fmt.Sprintf("Received an error from twiiter: %v", res.Errors))
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
		}`, true)

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twiiter: %v", res.Errors))
	}

	fmt.Println(res)
}