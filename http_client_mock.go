package twitterstream

import "net/http"

type mockHttpClient struct {
	token              string
	MockNewHttpRequest func(opts *requestOpts) (*http.Response, error)
	MockGetSearchStream func(queryParams string) (*http.Response, error)
	MockAddRules func(queryParams string, body string) (*http.Response, error)
	MockGenerateUrl func (name string, queryParams string) (string, error)
}

func newHttpClientMock(token string) *mockHttpClient {
	return &mockHttpClient{token: token}
}

func (t *mockHttpClient) generateUrl(name string, queryParams string) (string, error) {
	return t.MockGenerateUrl(name, queryParams)
}

func (t *mockHttpClient) addRules(queryParams string, body string) (*http.Response, error) {
	return t.MockAddRules(queryParams, body)
}

func (t *mockHttpClient) getSearchStream(queryParams string) (*http.Response, error) {
	return t.MockGetSearchStream(queryParams)
}

func (t *mockHttpClient) newHttpRequest(opts *requestOpts) (*http.Response, error) {
	return t.MockNewHttpRequest(opts)
}
