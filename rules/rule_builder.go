package rules

import "fmt"

type IRuleBuilder interface {
	Build() string
	WithValueAndTag(value string, tag string) string
	AddKeyword(keyword string) *RuleBuilder
	FilterKeyword(keyword string) *RuleBuilder
	AddExactPhrase(phrase string) *RuleBuilder
	FilterExactPhrase(phrase string) *RuleBuilder

}

type RuleBuilder struct {
	keywords []string
	filterkeywords []string
	exactphrase []string
	filterexactphrase []string
}

func NewRuleBuilder() IRuleBuilder {
	return &RuleBuilder{
		keywords: []string{},
		filterkeywords: []string{},
		exactphrase: []string{},
		filterexactphrase: []string{},
	}
}

func (r *RuleBuilder) Build() string {
	return "lol"
}

func (r *RuleBuilder) WithValueAndTag(value string, tag string) string {
	return fmt.Sprintf("%v %v", value, tag)
}

func (r *RuleBuilder) AddKeyword(keyword string) *RuleBuilder {
	r.keywords = append(r.keywords, keyword)
	return r
}

func (r *RuleBuilder) FilterKeyword(keyword string) *RuleBuilder {
	r.filterkeywords = append(r.filterkeywords, negate(keyword))
	return r
}

func (r *RuleBuilder) AddExactPhrase(phrase string) *RuleBuilder {
	r.exactphrase = append(r.exactphrase, phrase)
	return r
}

func (r *RuleBuilder) FilterExactPhrase(phrase string) *RuleBuilder {
	r.exactphrase = append(r.exactphrase, negate(phrase))
	return r
}

func negate(s string) string {
	ns := fmt.Sprintf("-%v", s)
	return ns
}

