package twitstream

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

type ITokenGenerator interface {
	RequestBearerToken() *requestBearerTokenResponse
	SetApiKeyAndSecret(apiKey, apiSecret string) *tokenGenerator
}

type tokenGenerator struct {
	apiKey    string
	apiSecret string
}

type requestBearerTokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

func newTokenGenerator() *tokenGenerator {
	return &tokenGenerator{}
}

func (a *tokenGenerator) SetApiKeyAndSecret(apiKey, apiSecret string) *tokenGenerator {
	a.apiKey = apiKey
	a.apiSecret = apiSecret
	return a
}

func (a *tokenGenerator) RequestBearerToken() *requestBearerTokenResponse {
	client := &http.Client{}
	body := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", body)

	if err != nil {
		log.Fatalf("Failed to construct request for bearer token: %v", err)
		return nil
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Authorization", "Basic "+a.base64EncodeKeys())

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed to perform request for bearer token: %v", err)
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
