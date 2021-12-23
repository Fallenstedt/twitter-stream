package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type twitterEndpoints map[string]string

// Endpoints is a map of twitter endpoints used to manage rules and streams.
var Endpoints = make(twitterEndpoints)

type (
	// IHttpClient is the interface the httpClient struct implements.
	IHttpClient interface {
		NewHttpRequest(opts *RequestOpts) (*http.Response, error)
		GetRules() (*http.Response, error)
		GetSearchStream(queryParams *url.Values) (*http.Response, error)
		AddRules(queryParams *url.Values, body string) (*http.Response, error)
		GenerateUrl(name string, queryParams *url.Values) (string, error)
	}

	httpClient struct {
		token string
	}
)

// NewHttpClient constructs a an HttpClient to interact with twitter.
func NewHttpClient(token string) IHttpClient {
	Endpoints["rules"] = "https://api.twitter.com/2/tweets/search/stream/rules"
	Endpoints["stream"] = "https://api.twitter.com/2/tweets/search/stream"
	Endpoints["token"] = "https://api.twitter.com/oauth2/token"
	return &httpClient{token}
}

// GetRules will return the current rules available for a specific API key.
func (t *httpClient) GetRules() (*http.Response, error) {
	res, err := t.NewHttpRequest(&RequestOpts{
		Method: "GET",
		Url:    Endpoints["rules"],
		Body:   "",
	})

	return res, err
}

// AddRules will add rules for you to stream with.
func (t *httpClient) AddRules(queryParams *url.Values, body string) (*http.Response, error) {
	url, err := t.GenerateUrl("rules", queryParams)

	if err != nil {
		return nil, err
	}

	res, err := t.NewHttpRequest(&RequestOpts{
		Method: "POST",
		Url:    url,
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetSearchStream will start the stream with twitter.
func (t *httpClient) GetSearchStream(queryParams *url.Values) (*http.Response, error) {
	// Make an HTTP GET request to GET /2/tweets/search/stream
	url, err := t.GenerateUrl("stream", queryParams)

	if err != nil {
		return nil, err
	}

	res, err := t.NewHttpRequest(&RequestOpts{
		Method: "GET",
		Url:    url,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GenerateUrl is a utility function for httpclient package to generate a valid url for api.twitter.
func (t *httpClient) GenerateUrl(name string, queryParams *url.Values) (string, error) {
	var url string
	if queryParams != nil {
		url = Endpoints[name] + fmt.Sprintf("?%v", queryParams.Encode())
	} else {
		url = Endpoints[name]
	}

	if len(url) == 0 || !strings.HasPrefix(url, "https://api.twitter.com") {
		return url, errors.New("Could not find endpoint with name " + name)
	} else {
		return url, nil
	}
}

// NewHttpRequest performs an authenticated http request with twitter with the token this httpclient has.
func (t *httpClient) NewHttpRequest(opts *RequestOpts) (*http.Response, error) {

	var req *http.Request
	var err error
	if opts.Method == "GET" {
		req, err = http.NewRequest(opts.Method, opts.Url, nil)
	} else {
		bufferBody := bytes.NewBuffer([]byte(opts.Body))
		req, err = http.NewRequest(opts.Method, opts.Url, bufferBody)
	}

	if err != nil {
		log.Printf("Failed to construct http request for %s: %v", opts.Url, err)
		return nil, err
	}

	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	if len(opts.Headers) > 0 {
		for _, header := range opts.Headers {
			req.Header.Set(header.Key, header.Value)
		}
	}

	// Set token if this httpclient has a token set
	if len(t.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}

	// Perform network request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to perform request for %s: %v", opts.Url, err)
		return nil, err
	}

	responseParser := new(httpResponseParser)
	return responseParser.handleResponse(resp, opts, t.NewHttpRequest)

}
