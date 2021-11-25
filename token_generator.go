package twitterstream

import (
	"encoding/base64"
	"encoding/json"
	"github.com/fallenstedt/twitter-stream/httpclient"
)

type (
	//ITokenGenerator is the interface that TokenGenerator implements.
	ITokenGenerator interface {
		RequestBearerToken() (*RequestBearerTokenResponse, error)
		SetApiKeyAndSecret(apiKey, apiSecret string) ITokenGenerator
	}
	TokenGenerator struct {
		httpClient httpclient.IHttpClient
		apiKey     string
		apiSecret  string
	}
	RequestBearerTokenResponse struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}
)

func newTokenGenerator(httpClient httpclient.IHttpClient) ITokenGenerator {
	return &TokenGenerator{httpClient: httpClient}
}

// SetApiKeyAndSecret sets the apiKey and apiSecret fields for the TokenGenerator instance.
func (a *TokenGenerator) SetApiKeyAndSecret(apiKey, apiSecret string) ITokenGenerator {
	a.apiKey = apiKey
	a.apiSecret = apiSecret
	return a
}

// RequestBearerToken requests a bearer token from twitter using the apiKey and apiSecret.
func (a *TokenGenerator) RequestBearerToken() (*RequestBearerTokenResponse, error) {

	resp, err := a.httpClient.NewHttpRequest(&httpclient.RequestOpts{
		Headers: []struct {
			Key   string
			Value string
		}{
			{"Content-Type", "application/x-www-form-urlencoded;charset=UTF-8"},
			{"Authorization", "Basic " + a.base64EncodeKeys()},
		},
		Method: "POST",
		Url:    httpclient.Endpoints["token"],
		Body:   "grant_type=client_credentials",
	})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data := new(RequestBearerTokenResponse)
	json.NewDecoder(resp.Body).Decode(data)

	return data, nil
}


func (a *TokenGenerator) base64EncodeKeys() string {
	// See Step 1 of encoding consumer key and secret twitter application-only requests here
	// https://developer.twitter.com/en/docs/authentication/oauth-2-0/application-only
	return base64.StdEncoding.EncodeToString([]byte(a.apiKey + ":" + a.apiSecret))
}
