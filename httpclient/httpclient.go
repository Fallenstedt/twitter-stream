package httpclient


import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

type twitterEndpoints map[string]string

var Endpoints = make(twitterEndpoints)

type (
	// IHttpClient is the interface the httpClient struct implements.
	IHttpClient interface {
		NewHttpRequest(opts *RequestOpts) (*http.Response, error)
		GetRules() (*http.Response, error)
		GetSearchStream(queryParams string) (*http.Response, error)
		AddRules(queryParams string, body string) (*http.Response, error)
		GenerateUrl(name string, queryParams string) (string, error)
	}

	httpClient struct {
		token string
	}

	RequestOpts struct {
		Retries uint8
		Method  string
		Url     string
		Body    string
		Headers []struct {
			Key   string
			Value string
		}
	}
)

func NewHttpClient(token string) *httpClient {
	Endpoints["rules"] = "https://api.twitter.com/2/tweets/search/stream/rules"
	Endpoints["stream"] = "https://api.twitter.com/2/tweets/search/stream"
	Endpoints["token"] = "https://api.twitter.com/oauth2/token"
	return &httpClient{token}
}

func (t *httpClient) GetRules() (*http.Response, error) {
	res, err := t.NewHttpRequest(&RequestOpts{
		Method: "GET",
		Url:    Endpoints["rules"],
		Body:   "",
	})

	return res, err
}

func (t *httpClient) AddRules(queryParams string, body string) (*http.Response, error) {
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

func (t *httpClient) GetSearchStream(queryParams string) (*http.Response, error) {
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

func (t *httpClient) GenerateUrl(name string, queryParams string) (string, error) {
	var url string
	if len(queryParams) > 0 {
		url = Endpoints[name] + queryParams
	} else {
		url = Endpoints[name]
	}

	if len(url) == 0 || !strings.HasPrefix(url, "https://api.twitter.com") {
		return url, errors.New("Could not find endpoint with name " + name)
	} else {
		return url, nil
	}
}

func (t *httpClient) NewHttpRequest(opts *RequestOpts) (*http.Response, error) {
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
			req.Header.Set(header.Key, header.Value)
		}
	}

	// Set token if this httpclient has a token set
	if len(t.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}

	// Perform network request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to perform request for %s: %v", opts.Url, err)
		return nil, err
	}

	// Retry with backoff if 429
	if resp.StatusCode == 429 {
		log.Printf("Retrying network request %s with backoff", opts.Url)

		delay := t.getBackOffTime(opts.Retries)
		log.Printf("Sleeping for %v seconds", delay)
		time.Sleep(delay)

		opts.Retries += 1
		return t.NewHttpRequest(opts)
	}

	// Reject if 400 or greater
	if resp.StatusCode >= 400 {
		log.Printf("Network Request at %s failed: %v", opts.Url, resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		msg := "Network request failed: " + string(body)
		return nil, errors.New(msg)
	}

	return resp, nil
}

func (t *httpClient) getBackOffTime(retries uint8) time.Duration {
	exponentialBackoffCeilingSecs := 30
	delaySecs := int(math.Floor((math.Pow(2, float64(retries)) - 1) * 0.5))
	if delaySecs > exponentialBackoffCeilingSecs {
		delaySecs = 30
	}
	return time.Duration(delaySecs) * time.Second
}
