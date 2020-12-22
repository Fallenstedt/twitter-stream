package twitstream

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

type (
	//ITokenGenerator is the interface that tokenGenerator implements
	ITokenGenerator interface {
		RequestBearerToken() *requestBearerTokenResponse
		SetApiKeyAndSecret(apiKey, apiSecret string) *tokenGenerator
	}
	tokenGenerator struct {
		httpClient IHttpClient
		apiKey     string
		apiSecret  string
	}
	requestBearerTokenResponse struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}
)

func newTokenGenerator(httpClient IHttpClient) *tokenGenerator {
	return &tokenGenerator{httpClient: httpClient}
}

// SetApiKeyAndSecret sets the apiLey and apiSecret fields for the instance
func (a *tokenGenerator) SetApiKeyAndSecret(apiKey, apiSecret string) *tokenGenerator {
	a.apiKey = apiKey
	a.apiSecret = apiSecret
	return a
}

// RequestBearerToken requests a bearer token from twitter using the apiKey and apiSecret
func (a *tokenGenerator) RequestBearerToken() *requestBearerTokenResponse {

	resp, err := a.httpClient.newHttpRequest(&requestOpts{
		Headers: []struct {
			key   string
			value string
		}{
			{"Content-Type", "application/x-www-form-urlencoded;charset=UTF-8"},
			{"Authorization", "Basic " + a.base64EncodeKeys()},
		},
		Method: "POST",
		Url:    endpoints["token"],
		Body:   "grant_type=client_credentials",
	})

	if err != nil {
		log.Fatalf("Failed to construct request for bearer token: %v", err)
		return nil
	}

	defer resp.Body.Close()
	data := new(requestBearerTokenResponse)
	json.NewDecoder(resp.Body).Decode(data)

	return data
}

func (a *tokenGenerator) base64EncodeKeys() string {
	return base64.StdEncoding.EncodeToString([]byte(a.apiKey + ":" + a.apiSecret))
}
