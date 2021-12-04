// Package twitterstream provides an easy way to stream tweets using Twitter's v2 Streaming API.
package twitterstream

import (
	"github.com/fallenstedt/twitter-stream/httpclient"
	"github.com/fallenstedt/twitter-stream/rules"
)

type TwitterApi struct {
	Rules  rules.IRules
	Stream IStream
}


// NewTokenGenerator creates a TokenGenerator which can request a Bearer token using a twitter api key and secret.
func NewTokenGenerator() ITokenGenerator {
	client := httpclient.NewHttpClient("")
	tokenGenerator := newTokenGenerator(client)
	return tokenGenerator
}

// NewTwitterStream consumes a twitter Bearer token.
// It is used to interact with Twitter's v2 filtered streaming API
func NewTwitterStream(token string) *TwitterApi {
	client := httpclient.NewHttpClient(token)
	rules := rules.NewRules(client)
	stream := newStream(client, newStreamResponseBodyReader())
	return &TwitterApi{Rules: rules, Stream: stream}
}
