package httpclient

import (
	"net/http"
	"testing"
)

func givenHttpResponseParserInstance() *httpResponseParser {
	return new(httpResponseParser)
}

func givenFakeHttpResponse(statusCode int) *http.Response {
	res := new(http.Response)
	res.StatusCode = statusCode
	return res
}

func TestHandleResponseShouldReturnIf200(t *testing.T) {
	instance := givenHttpResponseParserInstance()
	opts := new(RequestOpts)
	resp := givenFakeHttpResponse(200)

	result, err := instance.handleResponse(resp, opts, func(o *RequestOpts) (*http.Response, error) {
		return nil, nil
	})

	if err != nil {
		t.Errorf("Expected not error, got %v", err)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected a status code of 200")
	}
}

func TestHandleResponseShouldRetryRequestIf429(t *testing.T) {
	instance := givenHttpResponseParserInstance()
	opts := new(RequestOpts)
	resp := givenFakeHttpResponse(429)

	result, err := instance.handleResponse(resp, opts, func(o *RequestOpts) (*http.Response, error) {
		return givenFakeHttpResponse(200), nil
	})

	if opts.Retries != 1 {
		t.Errorf("Expected atleast on retry attempt, got %v", opts.Retries)
	}

	if err != nil {
		t.Errorf("Expected not error, got %v", err)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected a status code of 200")
	}
}

func TestHandleResponseShouldRejectIf400OrHigher(t *testing.T) {
	instance := givenHttpResponseParserInstance()
	opts := new(RequestOpts)
	resp := givenFakeHttpResponse(401)

	_, err := instance.handleResponse(resp, opts, func(o *RequestOpts) (*http.Response, error) {
		return nil, nil
	})

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
