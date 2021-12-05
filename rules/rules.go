package rules

import (
	"encoding/json"
	"github.com/fallenstedt/twitter-stream/httpclient"
)

type (
	//IRules is the interface the rules struct implements.
	IRules interface {
		AddRules(body string, dryRun bool) (*TwitterRuleResponse, error)
		Create(rules []*RuleValue, dryRun bool) (*TwitterRuleResponse, error)
		GetRules() (*TwitterRuleResponse, error)
	}

	//AddRulesRequest

	//TwitterRuleResponse is what is returned from twitter when adding or deleting a rule.
	TwitterRuleResponse struct {
		Data   []DataRule
		Meta   MetaRule
		Errors []ErrorRule
	}

	//DataRule is what is returned as "Data" when adding or deleting a rule.
	DataRule struct {
		Value string `json:"Value"`
		Tag   string `json:"Tag"`
		Id    string `json:"id"`
	}

	//MetaRule is what is returned as "Meta" when adding or deleting a rule.
	MetaRule struct {
		Sent    string      `json:"sent"`
		Summary MetaSummary `json:"summary"`
	}

	//MetaSummary is what is returned as "Summary" in "Meta" when adding or deleting a rule.
	MetaSummary struct {
		Created    uint `json:"created"`
		NotCreated uint `json:"not_created"`
	}

	//ErrorRule is what is returned as "Errors" when adding or deleting a rule.
	ErrorRule struct {
		Value string `json:"Value"`
		Id    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	}

	addRulesRequest struct {
		Add []*RuleValue `json:"add"`
	}
	rules struct {
		httpClient httpclient.IHttpClient
	}

)

//NewRules creates a "rules" instance. This is used to create Twitter Filtered Stream rules.
// https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/integrate/build-a-rule.
func NewRules(httpClient httpclient.IHttpClient) IRules {
	return &rules{httpClient: httpClient}
}

func (t *rules) Create(rules []*RuleValue, dryRun bool) (*TwitterRuleResponse, error) {
	add := addRulesRequest{Add: rules}
	body, err := json.Marshal(add)
	if err != nil {
		return nil, err
	}

	res, err := t.httpClient.AddRules(func() string {
		if dryRun {
			return "?dry_run=true"
		} else {
			return ""
		}
	}(), string(body))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data := new(TwitterRuleResponse)

	err = json.NewDecoder(res.Body).Decode(data)
	return data, err
}

// Deprecated: Use Create instead.
// AddRules adds or deletes rules to the stream using twitter's POST /2/tweets/search/stream/rules endpoint.
// The body is a stringified object.
// Learn about the possible error messages returned here https://developer.twitter.com/en/support/twitter-api/error-troubleshooting.
func (t *rules) AddRules(body string, dryRun bool) (*TwitterRuleResponse, error) {
	res, err := t.httpClient.AddRules(func() string {
		if dryRun {
			return "?dry_run=true"
		} else {
			return ""
		}
	}(), body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data := new(TwitterRuleResponse)

	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetRules gets rules for a stream using twitter's GET GET /2/tweets/search/stream/rules endpoint.
func (t *rules) GetRules() (*TwitterRuleResponse, error) {
	res, err := t.httpClient.GetRules()

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data := new(TwitterRuleResponse)
	json.NewDecoder(res.Body).Decode(data)

	return data, nil
}
