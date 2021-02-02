package twitterstream

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type twitterEndpoints map[string]string
var endpoints = make(twitterEndpoints)

type (
	// IHttpClient is the interface the httpClient struct implements.
	IHttpClient interface {

		newHttpRequest(opts *requestOpts) (*http.Response, error)
		getRules() (*http.Response, error)
		getSearchStream(queryParams string) (*http.Response, error)
		addRules(queryParams string, body string) (*http.Response, error)
		generateUrl(name string, queryParams string) (string, error)
	}

	httpClient struct {
		token string
	}

	requestOpts struct {
		Method  string
		Url     string
		Body    string
		Headers []struct {
			key   string
			value string
		}
	}
)

func newHttpClient(token string) *httpClient {
	endpoints["rules"] = "https://api.twitter.com/2/tweets/search/stream/rules"
	endpoints["stream"] = "https://api.twitter.com/2/tweets/search/stream"
	endpoints["token"] = "https://api.twitter.com/oauth2/token"
	return &httpClient{token}
}

func (t *httpClient) getRules() (*http.Response, error)  {
	res, err := t.newHttpRequest(&requestOpts{
		Method: "GET",
		Url:    endpoints["rules"],
		Body:   "",
	})

	return res, err
}

func (t *httpClient) addRules(queryParams string, body string) (*http.Response, error) {
	url, err :=  t.generateUrl("rules", queryParams)

	if err != nil {
		return nil, err
	}

	res, err := t.newHttpRequest(&requestOpts{
		Method: "POST",
		Url:    url,
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *httpClient) getSearchStream(queryParams string) (*http.Response, error) {
	// Make an HTTP GET request to GET /2/tweets/search/stream
	url, err := t.generateUrl("stream", queryParams)

	if err != nil {
		return nil, err
	}

	res, err := t.newHttpRequest(&requestOpts{
		Method: "GET",
		Url:    url,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *httpClient) generateUrl(name string, queryParams string) (string, error) {
	var url string
	if len(queryParams) > 0 {
		url = endpoints[name] + queryParams
	} else {
		url = endpoints[name]
	}

	if len(url) == 0 || !strings.HasPrefix(url, "https://api.twitter.com") {
		return url, errors.New("Could not find endpoint with name " + name)
	} else {
		return url, nil
	}
}

func (t *httpClient) newHttpRequest(opts *requestOpts) (*http.Response, error) {
	client := &http.Client{}

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
			req.Header.Set(header.key, header.value)
		}
	}

	// Set token if this client has a token set
	if len(t.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}

	// Perform network request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to perform request for %s: %v", opts.Url, err)
		return nil, err
	}

	// Reject if 400 or greater
	if resp.StatusCode >= 400 {
		log.Printf("Network Request failed: %v", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		msg := "Network request failed: " + string(body)
		return nil, errors.New(msg)
	}

	return resp, nil
}
