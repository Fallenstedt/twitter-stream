package httpclient

import "net/http"

type mockHttpClient struct {
	token               string
	MockNewHttpRequest  func(opts *RequestOpts) (*http.Response, error)
	MockGetSearchStream func(queryParams string) (*http.Response, error)
	MockGetRules        func() (*http.Response, error)
	MockAddRules        func(queryParams string, body string) (*http.Response, error)
	MockGenerateUrl     func(name string, queryParams string) (string, error)
}

func NewHttpClientMock(token string) *mockHttpClient {
	return &mockHttpClient{token: token}
}

func (t *mockHttpClient) GenerateUrl(name string, queryParams string) (string, error) {
	return t.MockGenerateUrl(name, queryParams)
}

func (t *mockHttpClient) GetRules() (*http.Response, error) {
	return t.MockGetRules()
}

func (t *mockHttpClient) AddRules(queryParams string, body string) (*http.Response, error) {
	return t.MockAddRules(queryParams, body)
}

func (t *mockHttpClient) GetSearchStream(queryParams string) (*http.Response, error) {
	return t.MockGetSearchStream(queryParams)
}

func (t *mockHttpClient) NewHttpRequest(opts *RequestOpts) (*http.Response, error) {
	return t.MockNewHttpRequest(opts)
}
