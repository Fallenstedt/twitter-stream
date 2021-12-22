package rules

import (
	"encoding/json"
	"github.com/fallenstedt/twitter-stream/httpclient"
	"net/url"
)

type (
	//IRules is the interface the rules struct implements.
	IRules interface {
		Create(rules CreateRulesRequest, dryRun bool) (*TwitterRuleResponse, error)
		Delete(req DeleteRulesRequest, dryRun bool) (*TwitterRuleResponse, error)
		Get() (*TwitterRuleResponse, error)
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

	rules struct {
		httpClient httpclient.IHttpClient
	}

)

//NewRules creates a "rules" instance. This is used to create Twitter Filtered Stream rules.
// https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/integrate/build-a-rule.
func NewRules(httpClient httpclient.IHttpClient) IRules {
	return &rules{httpClient: httpClient}
}

// Create will create new twitter streaming rules.
func (t *rules) Create(rules CreateRulesRequest, dryRun bool) (*TwitterRuleResponse, error) {
	body, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}

	res, err := t.httpClient.AddRules(t.addDryRun(dryRun), string(body))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data := new(TwitterRuleResponse)

	err = json.NewDecoder(res.Body).Decode(data)
	return data, err
}
// Delete will delete rules twitter rules by their id.
func (t *rules) Delete(req DeleteRulesRequest, dryRun bool) (*TwitterRuleResponse, error) {

	body, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := t.httpClient.AddRules(t.addDryRun(dryRun), string(body))


	defer res.Body.Close()
	data := new(TwitterRuleResponse)

	err = json.NewDecoder(res.Body).Decode(data)
	return data, err
}


// Get will fetch the current rules.
func (t *rules) Get() (*TwitterRuleResponse, error) {
	res, err := t.httpClient.GetRules()

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data := new(TwitterRuleResponse)
	json.NewDecoder(res.Body).Decode(data)

	return data, nil
}



func (t *rules) addDryRun(dryRun bool) *url.Values {
	if dryRun {
		query := new(url.URL).Query()
		query.Add("dry_run", "true")
		return &query
	} else {
		return nil
	}
}