package twitstream

import "net/http"

type mockHttpClient struct {
	token              string
	MockNewHttpRequest func(opts *requestOpts) (*http.Response, error)
}

func newHttpClientMock(token string) *mockHttpClient {
	return &mockHttpClient{token: token}
}

func (t *mockHttpClient) newHttpRequest(opts *requestOpts) (*http.Response, error) {
	return t.MockNewHttpRequest(opts)
}