package twitstream

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSetApiKeyAndSecret(t *testing.T) {
	var tests = []struct {
		apiKey    string
		apiSecret string
		result    tokenGenerator
	}{
		{"foo", "bar", tokenGenerator{apiKey: "foo", apiSecret: "bar"}},
		{"", "", tokenGenerator{apiKey: "", apiSecret: ""}},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("(%d) %s %s", i, tt.apiKey, tt.apiSecret)
		t.Run(testName, func(t *testing.T) {
			result := newTokenGenerator(newHttpClientMock(""))
			result.SetApiKeyAndSecret(tt.apiKey, tt.apiSecret)

			if result.apiKey != tt.result.apiKey {
				t.Errorf("got %s, want %s", result.apiKey, tt.result.apiKey)
			}

			if result.apiKey != tt.result.apiKey {
				t.Errorf("got %s, want %s", result.apiKey, tt.result.apiKey)
			}
		})
	}
}

func TestRequestBearerToken(t *testing.T) {
	var tests = []struct {
		mockRequest func(opts *requestOpts) (*http.Response, error)
		result      *requestBearerTokenResponse
	}{
		{func(opts *requestOpts) (*http.Response, error) {

			json := `{
				"token_type": "bearer",
				"access_token": "123Token456"
			}`
			// create a new reader with that JSON
			body := ioutil.NopCloser(bytes.NewReader([]byte(json)))

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       body,
			}, nil
		},
			&requestBearerTokenResponse{
				TokenType:   "bearer",
				AccessToken: "123Token456",
			}},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("(%d)", i)

		t.Run(testName, func(t *testing.T) {
			mockClient := newHttpClientMock("")
			mockClient.MockNewHttpRequest = tt.mockRequest

			instance := newTokenGenerator(mockClient)
			instance.SetApiKeyAndSecret("SomeKey", "SomeSecret")

			data := instance.RequestBearerToken()

			if data == nil {
				t.Errorf("got %s, want %s", data, tt.result)
			}

			if data.AccessToken != tt.result.AccessToken {
				t.Errorf("got %s, want %s", data.AccessToken, tt.result.AccessToken)
			}

			if data.TokenType != tt.result.TokenType {
				t.Errorf("got %s, want %s", data.TokenType, tt.result.TokenType)
			}
		})
	}
}
