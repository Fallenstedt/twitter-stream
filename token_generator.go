package twitterstream

import (
	"encoding/base64"
	"encoding/json"
	"github.com/fallenstedt/twitter-stream/httpclient"
)

type (
	//ITokenGenerator is the interface that tokenGenerator implements.
	ITokenGenerator interface {
		RequestBearerToken() *requestBearerTokenResponse
		SetApiKeyAndSecret(apiKey, apiSecret string) *tokenGenerator
	}
	tokenGenerator struct {
		httpClient httpclient.IHttpClient
		apiKey     string
		apiSecret  string
	}
	requestBearerTokenResponse struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}
)

func newTokenGenerator(httpClient httpclient.IHttpClient) *tokenGenerator {
	return &tokenGenerator{httpClient: httpClient}
}

// SetApiKeyAndSecret sets the apiKey and apiSecret fields for the tokenGenerator instance.
func (a *tokenGenerator) SetApiKeyAndSecret(apiKey, apiSecret string) *tokenGenerator {
	a.apiKey = apiKey
	a.apiSecret = apiSecret
	return a
}

// RequestBearerToken requests a bearer token from twitter using the apiKey and apiSecret.
func (a *tokenGenerator) RequestBearerToken() (*requestBearerTokenResponse, error) {

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
	data := new(requestBearerTokenResponse)
	json.NewDecoder(resp.Body).Decode(data)

	return data, nil
}

func (a *tokenGenerator) base64EncodeKeys() string {
	return base64.StdEncoding.EncodeToString([]byte(a.apiKey + ":" + a.apiSecret))
}
