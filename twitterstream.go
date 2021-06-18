// Package twitterstream provides an easy way to stream tweets using Twitter's v2 Streaming API.
package twitterstream

import "github.com/fallenstedt/twitter-stream/httpclient"

type twitterApi struct {
	Rules  IRules
	Stream IStream
}

// NewTokenGenerator creates a tokenGenerator which can request a Bearer token using a twitter api key and secret.
func NewTokenGenerator() *tokenGenerator {
	client := httpclient.NewHttpClient("")
	tokenGenerator := newTokenGenerator(client)
	return tokenGenerator
}

// NewTwitterStream consumes a twitter Bearer token.
// It is used to interact with Twitter's v2 filtered streaming API
func NewTwitterStream(token string) *twitterApi {
	client := httpclient.NewHttpClient(token)
	rules := newRules(client)
	stream := newStream(client, newStreamResponseBodyReader())
	return &twitterApi{Rules: rules, Stream: stream}
}
