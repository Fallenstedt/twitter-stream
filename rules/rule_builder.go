package rules

type  (
	// IRuleBuilder is an interface that describers how to implement a RuleBuilder.
	IRuleBuilder interface {
		AddRule(value string, tag string) *RuleBuilder
		Build() CreateRulesRequest
	}

	// RuleValue is a struct used to help create twitter stream rules.
	// It takes in a value and a tag.
	RuleValue struct {
		Value *string `json:"value,omitempty"`
		Tag   *string `json:"tag,omitempty"`
	}

	// RuleBuilder is struct used to help create twitter stream rules.
	RuleBuilder struct {
		rules []*RuleValue
	}

	// CreateRulesRequest is a struct used to create the payload for creating rules.
	CreateRulesRequest struct {
		Add []*RuleValue `json:"add"`
	}

	// DeleteRulesRequest is a struct used to create the payload for deleting rules.
	DeleteRulesRequest struct {
		Delete struct {
			Ids []int `json:"ids"`
		} `json:"delete"`
	}

)

// NewDeleteRulesRequest will create an instance of DeleteRulesRequest.
func NewDeleteRulesRequest(ids ...int) DeleteRulesRequest {
	return DeleteRulesRequest{Delete: struct {
		Ids []int `json:"ids"`
	}(struct{ Ids []int }{Ids: ids})}
}

// NewRuleBuilder will create an instance of `RuleBuilder`.
func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{
		rules: []*RuleValue{},
	}
}

//AddRule will create a rule to be build for filtered-stream.
//Read more about rule limitations here https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction.
func (r *RuleBuilder) AddRule(value string, tag string) *RuleBuilder {
	rule := newRuleValue().setValueTag(value, tag)
	r.rules = append(r.rules, rule)
	return r
}

func (r *RuleBuilder) Build() CreateRulesRequest {
	add := CreateRulesRequest{Add: r.rules}
	return add
}

func newRuleValue() *RuleValue {
	return &RuleValue{
		Value: nil,
		Tag:   nil,
	}
}

func (r *RuleValue) setValueTag(value string, tag string) *RuleValue {
	r.Value = &value
	r.Tag = &tag
	return r
}
