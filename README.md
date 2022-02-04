# TwitterStream

![go twitter](./go-twitter.png)

[![v2](https://img.shields.io/endpoint?url=https%3A%2F%2Ftwbadges.glitch.me%2Fbadges%2Fv2)](https://developer.twitter.com/en/docs/twitter-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/fallenstedt/twitter-stream)](https://goreportcard.com/report/github.com/fallenstedt/twitter-stream)
[![Go Reference](https://pkg.go.dev/badge/github.com/fallenstedt/twitter-stream.svg)](https://pkg.go.dev/github.com/fallenstedt/twitter-stream)

TwitterStream is a Go library for creating streaming rules and streaming tweets with [Twitter's v2 Filtered Streaming API](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction).

- See [my blog post](https://www.fallenstedt.com/blog/twitter-stream/) for a tutorial on Twitter's Filtered Stream endpoint.
- See [examples](https://github.com/fallenstedt/twitter-stream/tree/master/example) to start adding your own rules and start streaming.

![example of twit stream](./example.gif)

## Installation

`go get github.com/fallenstedt/twitter-stream`


## Examples

See [examples](https://github.com/fallenstedt/twitter-stream/tree/master/example), or follow the guide below.

#### Starting a stream

##### Obtain an Access Token using your Twitter Access Key and Secret.

You need an access token to do any streaming. `twitterstream` provides an easy way to fetch an access token. Use your
API key and secret API key from twitter to request an access token.

```go
	tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret("key", "secret").RequestBearerToken()
```

##### Create a streaming api

Create a twitterstream instance with your access token from above.

```go
	api := twitterstream.NewTwitterStream(tok.AccessToken)
```

##### Create rules

We need to create [twitter streaming rules](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/integrate/build-a-rule) so we can get tweets that we want.
The filtered stream endpoints deliver filtered Tweets to you in real-time that match on a set of rules that are applied to the stream. Rules are made up of operators that are used to match on a variety of Tweet attributes.
Below we create three rules. One for puppy tweets with images, another for cat tweets with images, and the other of unique English golang job postings. Each rule is
associated with their own tag.

```go

rules := twitterstream.NewRuleBuilder().
            AddRule("cat has:images", "cat tweets with images").
            AddRule("puppy has:images", "puppy tweets with images").
            AddRule("lang:en -is:retweet -is:quote (#golangjobs OR #gojobs)", "golang jobs").
            Build()

// Create will create twitter rules
// dryRun is set to false. Set to true to test out your request
res, err := api.Rules.Create(rules, false)

// Get will get your current rules
res, err := api.Rules.Get()

// Delete will delete your rules by their id
// dryRun is set to false. Set to true to test out your request
res, err := api.Rules.Delete(rules.NewDeleteRulesRequest(1468427075727945728, 1468427075727945729), false)


```

##### Set your unmarshal hook

It is encouraged you set an unmarshal hook for thread-safety. Go's `bytes.Buffer` is not thread safe. Sharing a `bytes.Buffer`
across multiple goroutines introduces risk of panics when decoding json.
To avoid panics, it's encouraged to unmarshal json in the same goroutine where the `bytes.Buffer` exists. Use `SetUnmarshalHook` to set a function that unmarshals json.

By default, twitterstream's unmarshal hook will return `[]byte` if you want to live dangerously.

```go

type StreamDataExample struct {
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
        ID  string `json:"id"`
        Tag string `json:"tag"`
    } `json:"matching_rules"`
}

api.SetUnmarshalHook(func(bytes []byte) (interface{}, error) {
    data := StreamDataExample{}

    if err := json.Unmarshal(bytes, &data); err != nil {
        fmt.Printf("failed to unmarshal bytes: %v", err)
    }

    return data, err
})
```

##### Start Stream

Start your stream. This is a long-running HTTP GET request.
You can request additional tweet data by adding [query params](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream).
Use the `twitterstream.NewStreamQueryParamsBuilder()` to start a stream with the data you want.

```go

// Steps from above, Placed into a single function
// This assumes you have at least one streaming rule configured.
// returns a configured instance of twitterstream
func fetchTweets() stream.IStream {
    tok, err := twitterstream.NewTokenGenerator().SetApiKeyAndSecret(KEY, SECRET).RequestBearerToken()

    if err != nil {
        panic(err)
    }

    api := twitterstream.NewTwitterStream(tok).Stream
    api.SetUnmarshalHook(func(bytes []byte) (interface{}, error) {
        data := StreamDataExample{}

        if err := json.Unmarshal(bytes, &data); err != nil {
          fmt.Printf("failed to unmarshal bytes: %v", err)
        }
        return data, err
    })

    // https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream
    streamExpansions := twitterstream.NewStreamQueryParamsBuilder().
        AddExpansion("author_id").
        AddTweetField("created_at").
        Build()

    // StartStream will start the stream
    err = api.StartStream(streamExpansions)

    if err != nil {
        panic(err)
    }

    return api
}

// This will run forever
func initiateStream() {
    fmt.Println("Starting Stream")

    // Start the stream
    // And return the library's api
    api := fetchTweets()

    // When the loop below ends, restart the stream
    defer initiateStream()

    // Start processing data from twitter after starting the stream
    for tweet := range api.GetMessages() {

        // Handle disconnections from twitter
        // https://developer.twitter.com/en/docs/twitter-api/tweets/volume-streams/integrate/handling-disconnections
        if tweet.Err != nil {
            fmt.Printf("got error from twitter: %v", tweet.Err)

            // Notice we "StopStream" and then "continue" the loop instead of breaking.
            // StopStream will close the long running GET request to Twitter's v2 Streaming endpoint by
            // closing the `GetMessages` channel. Once it's closed, it's safe to perform a new network request
            // with `StartStream`
            api.StopStream()
            continue
        }
        result := tweet.Data.(StreamDataExample)

        // Here I am printing out the text.
        // You can send this off to a queue for processing.
        // Or do your processing here in the loop
        fmt.Println(result.Data.Text)
    }

    fmt.Println("Stopped Stream")
}
```

## Contributing

Pull requests and feature requests are always welcome.
Please accompany a pull request with tests.
