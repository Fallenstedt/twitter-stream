package twitterstream

import (
	"encoding/json"
)

type (
	//IRules is the interface the rules struct implements.
	IRules interface {
		AddRules(body string, dryRun bool) (*rulesResponse, error)
		GetRules() (*rulesResponse, error)
	}

	rules struct {
		httpClient IHttpClient
	}

	rulesResponse struct {
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

// AddRules adds or deletes rules to the stream using twitter's POST /2/tweets/search/stream/rules endpoint.
// The body is a stringified object.
func (t *rules) AddRules(body string, dryRun bool) (*rulesResponse, error) {

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
	data := new(rulesResponse)
	json.NewDecoder(res.Body).Decode(data)

	return data, nil
}

// GetRules gets rules for a stream using twitter's GET GET /2/tweets/search/stream/rules endpoint.
func (t *rules) GetRules() (*rulesResponse, error) {
	res, err := t.httpClient.newHttpRequest(&requestOpts{
		Method: "GET",
		Url:    endpoints["rules"],
		Body:   "",
	})

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data := new(rulesResponse)
	json.NewDecoder(res.Body).Decode(data)

	return data, nil
}
