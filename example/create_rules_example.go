package main

import (
	"fmt"
	twitterstream "github.com/fallenstedt/twitter-stream"
	"github.com/fallenstedt/twitter-stream/rules"
)

func addRules() {

	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(KEY, SECRET).RequestBearerToken()
	if err != nil {
		panic(err)
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken)
	rules := twitterstream.NewRuleBuilder().
		AddRule("cat has:images", "cat tweets with images").
		AddRule("puppy has:images", "puppy tweets with images").
		AddRule("lang:en -is:retweet -is:quote (#golangjobs OR #gojobs)", "golang jobs").
		Build()

	res, err := api.Rules.Create(rules, false) // dryRun is set to false.

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twitter: %v", res.Errors))
	}

	fmt.Println("I have created rules.")
	printRules(res.Data)
}

func getRules() {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(KEY, SECRET).RequestBearerToken()
	if err != nil {
		panic(err)
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken)
	res, err := api.Rules.Get()

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twitter: %v", res.Errors))
	}

	if len(res.Data) > 0 {
		fmt.Println("I found these rules: ")
		printRules(res.Data)
	} else {
		fmt.Println("I found no rules")
	}

}

func deleteRules() {
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(KEY, SECRET).RequestBearerToken()
	if err != nil {
		panic(err)
	}
	api := twitterstream.NewTwitterStream(tok.AccessToken)

	// use api.Rules.Get to find the ID number for an existing rule
	res, err := api.Rules.Delete(rules.NewDeleteRulesRequest(1469776000158363653, 1469776000158363654), false)

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twitter: %v", res.Errors))
	}

	fmt.Println("I have deleted rules ")
}


func printRules(data []rules.DataRule) {
	for _, datum := range data {
		fmt.Printf("Id: %v\n", datum.Id)
		fmt.Printf("Tag: %v\n",datum.Tag)
		fmt.Printf("Value: %v\n\n", datum.Value)
	}
}