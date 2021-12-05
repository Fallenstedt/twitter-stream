package main

import (
	"fmt"
	twitterstream "github.com/fallenstedt/twitter-stream"
)

const key = "YOUR_KEY"
const secret = "YOUR_SECRET"

func main() {
 	addRules()
 	//getRules()
 	//deleteRules()
}

func addRules() {

	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(key, secret).RequestBearerToken()
	if err != nil {
		panic(err)
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken)
	rules := twitterstream.NewRuleBuilder().
		AddRule("puppies has:images", "puppy tweets with images").
		AddRule("lang:en -is:retweet -is:quote (#golangjobs OR #gojobs)", "golang jobs").
		Build()

	res, err := api.Rules.Create(rules, true) // dryRun is set to false.

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
