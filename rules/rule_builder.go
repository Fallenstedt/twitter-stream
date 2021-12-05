package rules

type  (
	IRuleBuilder interface {
		AddRule(value string, tag string) *RuleBuilder
		Build() []*RuleValue
	}

	RuleValue struct {
		Value *string `json:"value,omitempty"`
		Tag   *string `json:"tag,omitempty"`
	}

	RuleBuilder struct {
		rules []*RuleValue
	}
)

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

func (r *RuleBuilder) Build() []*RuleValue {
	return r.rules
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

