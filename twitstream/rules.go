package twitstream

import (
	"encoding/json"
	"fmt"
)

type (
	//IRules is the interface the rules struct implements
	IRules interface {
		AddRules(body string, dryRun bool) (*addRulesResponse, error)
		GetRules() (*getRulesResponse, error)
	}

	rules struct {
		httpClient IHttpClient
	}

	getRulesResponse struct {
		Data []rulesResponseValue
		Meta rulesResponseMeta
	}

	addRulesResponse struct {
		Data []rulesResponseValue
		Meta rulesResponseMeta
	}

	rulesResponseValue struct {
		Value string `json:"value"`
		Tag   string `json:"tag"`
		Id    string `json:"id"`
	}
	rulesResponseMeta struct {
		Sent    string                      `json:"sent"`
		Summary addRulesResponseMetaSummary `json:"summary"`
	}
	addRulesResponseMetaSummary struct {
		Created    uint `json:"created"`
		NotCreated uint `json:"not_created"`
	}
)

func newRules(httpClient IHttpClient) *rules {
	return &rules{httpClient: httpClient}
}

// AddRules adds rules to the stream. body is stringified object:
// https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/post-tweets-search-stream-rules
func (t *rules) AddRules(body string, dryRun bool) (*addRulesResponse, error) {

	var url string
	if dryRun {
		url = endpoints["rules"] + "?dry_run=true"
	} else {
		url = endpoints["rules"]
	}

	res, err := t.httpClient.newHttpRequest(&requestOpts{
		Method: "POST",
		Url:    url,
		Body:   body,
	})

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data := new(addRulesResponse)
	json.NewDecoder(res.Body).Decode(data)

	return data, nil
}

// GetRules gets rules for a stream.
// https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/get-tweets-search-stream-rules
func (t *rules) GetRules() (*getRulesResponse, error) {
	res, err := t.httpClient.newHttpRequest(&requestOpts{
		Method: "GET",
		Url:    endpoints["rules"],
		Body:   "",
	})

	if err != nil {
		return nil, err
	}
fmt.Println("test")
	defer res.Body.Close()
	data := new(getRulesResponse)
	json.NewDecoder(res.Body).Decode(data)

	return data, nil
}
