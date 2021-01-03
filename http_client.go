package twitter_stream

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

var endpoints = make(map[string]string)

type (
	// IHttpClient is the interface the httpClient struct implements.
	IHttpClient interface {
		newHttpRequest(opts *requestOpts) (*http.Response, error)
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
