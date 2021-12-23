package rules

import (
	"encoding/json"
	"testing"
)

func TestNewDeleteRulesRequestAcceptsMultipleIds(t *testing.T) {
	result := NewDeleteRulesRequest(1, 2, 3, 4, 5)
	expected := []int{1, 2, 3, 4, 5}

	if len(result.Delete.Ids) != len(expected) {
		t.Errorf("Expected %v to have same length as %v", result, expected)
	}

	for i := range result.Delete.Ids {
		if result.Delete.Ids[i] != expected[i] {
			t.Errorf("Expected %v to equal %v", result.Delete.Ids[i], expected[i])
		}
	}
}

func TestNewDeleteRulesRequestMarshalsWell(t *testing.T) {
	result := NewDeleteRulesRequest(1, 2, 3, 4, 5)
	body, err := json.Marshal(result)

	if err != nil {
		t.Error(err)
	}

	if string(body) != "{\"delete\":{\"ids\":[1,2,3,4,5]}}" {
		t.Errorf("Expected %v to equal %v", string(body), "{\"delete\":{\"ids\":[1,2,3,4,5]}}")
	}
}

func TestNewRuleBuilderBuildsManyRules(t *testing.T) {
	result := NewRuleBuilder().AddRule("cats", "cat tweets").AddRule("dogs", "dog tweets").Build()
	body, err := json.Marshal(result)

	if err != nil {
		t.Error(err)
	}

	if string(body) != "{\"add\":[{\"value\":\"cats\",\"tag\":\"cat tweets\"},{\"value\":\"dogs\",\"tag\":\"dog tweets\"}]}" {
		t.Errorf("Expected %v to equal %v", string(body), "{\"add\":[{\"value\":\"cats\",\"tag\":\"cat tweets\"},{\"value\":\"dogs\",\"tag\":\"dog tweets\"}]}")
	}
}
