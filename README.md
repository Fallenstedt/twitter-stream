# TwitterStream

![go twitter](./go-twitter.png)

[![v2](https://img.shields.io/endpoint?url=https%3A%2F%2Ftwbadges.glitch.me%2Fbadges%2Fv2)](https://developer.twitter.com/en/docs/twitter-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/fallenstedt/twitter-stream)](https://goreportcard.com/report/github.com/fallenstedt/twitter-stream)
[![Go Reference](https://pkg.go.dev/badge/github.com/fallenstedt/twitter-stream.svg)](https://pkg.go.dev/github.com/fallenstedt/twitter-stream)

TwitStream is a Go library for creating streaming rules and streaming tweets with [Twitter's v2 Filtered Streaming API](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction). 
See [examples](https://github.com/fallenstedt/twitter-stream/tree/master/example) to start adding your own rules and start streaming.  


![example of twit stream](./example.gif)


## Installation

`go get github.com/fallenstedt/twitter-stream`




## Examples

#### Starting a stream

##### Obtain an Access Token using your Twitter Access Key and Secret.
You need an access token to do any streaming. `twitterstream` provides an easy way to fetch an access token.
```go
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret("key", "secret").RequestBearerToken()

	if err != nil {
		panic(err)
	}
```

##### Create a streaming api
Create a twitterstream instance with your access token from above.

```go
	api := twitterstream.NewTwitterStream(tok.AccessToken)
```

##### Start Stream
Start your stream. This is a long-running HTTP GET request. 
You can get specific data you want by adding [query params](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream).
Additionally, [view an example of query params here](https://developer.twitter.com/en/docs/twitter-api/expansions).

```go
    err := api.Stream.StartStream("")

	if err != nil {
		panic(err)
	}
```

4. Consume Messages from the Stream
Handle any `io.EOF` and other errors that arise first, then unmarshal your bytes into your favorite struct. Below is an example with strings 
```go
	go func() {
		for message := range api.Stream.GetMessages() {
			if message.Err != nil {
				panic(message.Err)
			}
                        // Will print something like: 
                        //{"data":{"id":"1356479201000","text":"Look at this cat picture"},"matching_rules":[{"id":12345,"tag":"cat tweets with images"}]}
			fmt.Println(string(message.Data))
		}
	}()

	time.Sleep(time.Second * 30)
	api.Stream.StopStream()
```

#### Creating, Deleting, and Getting Rules

##### Obtain an Access Token using your Twitter Access Key and Secret.
You need an access token to do anything. `twitterstream` provides an easy way to fetch an access token.
```go
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret("key", "secret").RequestBearerToken()

	if err != nil {
		panic(err)
	}
```

##### Create a streaming api
Create a twitterstream instance with your access token from above.

```go
	api := twitterstream.NewTwitterStream(tok.AccessToken)
```

##### Get Rules
Use the `Rules` struct to access different Rules endpoints as defined in [Twitter's API Reference](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference)
```go
    res, err := api.Rules.GetRules()

	if err != nil {
		panic(err)
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		//https://developer.twitter.com/en/support/twitter-api/error-troubleshooting
		panic(fmt.Sprintf("Received an error from twiiter: %v", res.Errors))
	}

	fmt.Println(res.Data)
```

##### Add Rules
```go
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
```
##### Delete Rules
```go
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

```


## Contributing

Pull requests are always welcome. Please accompany a pull request with tests. 

