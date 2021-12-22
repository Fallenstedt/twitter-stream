package httpclient

import (
	"net/http"
	"net/url"
)

type mockHttpClient struct {
	token               string
	MockNewHttpRequest  func(opts *RequestOpts) (*http.Response, error)
	MockGetSearchStream func(queryParams *url.Values) (*http.Response, error)
	MockGetRules        func() (*http.Response, error)
	MockAddRules        func(queryParams *url.Values, body string) (*http.Response, error)
	MockGenerateUrl     func(name string, queryParams *url.Values) (string, error)
}

func NewHttpClientMock(token string) *mockHttpClient {
	return &mockHttpClient{token: token}
}

func (t *mockHttpClient) GenerateUrl(name string, queryParams *url.Values) (string, error) {
	return t.MockGenerateUrl(name, queryParams)
}

func (t *mockHttpClient) GetRules() (*http.Response, error) {
	return t.MockGetRules()
}

func (t *mockHttpClient) AddRules(queryParams *url.Values, body string) (*http.Response, error) {
	return t.MockAddRules(queryParams, body)
}

func (t *mockHttpClient) GetSearchStream(queryParams *url.Values) (*http.Response, error) {
	return t.MockGetSearchStream(queryParams)
}

func (t *mockHttpClient) NewHttpRequest(opts *RequestOpts) (*http.Response, error) {
	return t.MockNewHttpRequest(opts)
}
