package main

import (
	"fmt"
	twitterstream "github.com/fallenstedt/twitter-stream"
	"time"
)

const key = "YOUR_KEY"
const secret = "YOUR_SECRET"


func main() {
	// Use your favorite function from below here
 	startStream()
 	//addRules()
 	//getRules()
	//deleteRules()
}


func startStream() {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(key, secret).RequestBearerToken()

	if err != nil {
		panic(err)
	}

	api := twitterstream.NewTwitterStream(tok.AccessToken)

	err = api.Stream.StartStream("")

	if err != nil {
		panic(err)
	}

	go func() {
		for message := range api.Stream.GetMessages() {
			if message.Err != nil {
				panic(message.Err)
			}
			fmt.Println(string(message.Data))
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