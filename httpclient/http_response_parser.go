package httpclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"
)

// httpResponseParser is a struct that will retry network requests if the response has a status code of 429.
type httpResponseParser struct{}

func (h httpResponseParser) handleResponse(resp *http.Response, opts *RequestOpts, fn func(opts *RequestOpts) (*http.Response, error)) (*http.Response, error) {
	// Retry with backoff if 429
	if resp.StatusCode == 429 {
		log.Printf("Retrying network request %s with backoff", opts.Url)

		delay := h.getBackOffTime(opts.Retries)
		log.Printf("Sleeping for %v seconds", delay)
		time.Sleep(delay)

		opts.Retries += 1

		return fn(opts)
	}

	// Reject if 400 or greater
	if resp.StatusCode >= 400 {
		log.Printf("Network Request at %s failed: %v", opts.Url, resp.StatusCode)

		var msg string
		if resp.Body != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			msg = "Network request failed: " + string(body)
		} else {
			msg = "Network request failed with status" + fmt.Sprint(resp.StatusCode)
		}

		return nil, errors.New(msg)
	}

	return resp, nil
}

func (h httpResponseParser) getBackOffTime(retries uint8) time.Duration {
	exponentialBackoffCeilingSecs := 30
	delaySecs := int(math.Floor((math.Pow(2, float64(retries)) - 1) * 0.5))
	if delaySecs > exponentialBackoffCeilingSecs {
		delaySecs = 30
	}
	return time.Duration(delaySecs) * time.Second
}
