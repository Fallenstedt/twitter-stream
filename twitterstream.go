// Package twitterstream provides an easy way to stream tweets using Twitter's v2 Streaming API.
package twitterstream

import (
	"github.com/fallenstedt/twitter-stream/httpclient"
	"github.com/fallenstedt/twitter-stream/rules"
	"github.com/fallenstedt/twitter-stream/stream"
	"github.com/fallenstedt/twitter-stream/token_generator"
)

type TwitterApi struct {
	Rules  rules.IRules
	Stream stream.IStream
}

// NewTokenGenerator creates a TokenGenerator which can request a Bearer token using a twitter api key and secret.
func NewTokenGenerator() token_generator.ITokenGenerator {
	client := httpclient.NewHttpClient("")
	tokenGenerator := token_generator.NewTokenGenerator(client)
	return tokenGenerator
}

func NewRuleBuilder() rules.IRuleBuilder {
	return rules.NewRuleBuilder()
}

// NewTwitterStream consumes a twitter Bearer token.
// It is used to interact with Twitter's v2 filtered streaming API
func NewTwitterStream(token string) *TwitterApi {
	client := httpclient.NewHttpClient(token)
	rules := rules.NewRules(client)
	stream := stream.NewStream(client, stream.NewStreamResponseBodyReader())
	return &TwitterApi{Rules: rules, Stream: stream}
}
