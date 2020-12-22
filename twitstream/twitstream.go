package twitstream

type twitterApi struct {
	Rules  IRules
	Stream IStream
}

func NewTokenGenerator() *tokenGenerator {
	client := newHttpClient("")
	tokenGenerator := newTokenGenerator(client)
	return tokenGenerator
}

func NewTwitterStream(token string) *twitterApi {
	client := newHttpClient(token)
	rules := newRules(client)
	stream := newStream(client, newStreamResponseBodyReader())
	return &twitterApi{Rules: rules, Stream: stream}
}
