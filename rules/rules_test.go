package rules

import (
	"bytes"
	"fmt"
	"github.com/fallenstedt/twitter-stream/httpclient"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {

	var tests = []struct {
		body        CreateRulesRequest
		mockRequest func(queryParams string, body string) (*http.Response, error)
		result      *TwitterRuleResponse
	}{
		{
			NewRuleBuilder().AddRule("cat has:images", "cat tweets with images").Build(),
			func(queryParams string, bodyRequest string) (*http.Response, error) {
				json := `{
					"data": [{
						"Value": "cat has:images", 
						"Tag":"cat tweets with images", 
						"id": "123456"
					}],
					"meta": {
						"sent": "today",
						"summary": {
							"created": 1,
							"not_created": 0
						}
					},
					"errors": []
				}`

				body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
				}, nil

			},
			&TwitterRuleResponse{
				Data: []DataRule{
					{
						Value: "cat has:images",
						Tag:   "cat tweets with images",
						Id:    "123456",
					},
				},
				Meta: MetaRule{
					Sent: "today",
					Summary: MetaSummary{
						Created:    1,
						NotCreated: 0,
					},
				},
				Errors: nil,
			},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("TestCreate (%d)", i)

		t.Run(testName, func(t *testing.T) {
			mockClient := httpclient.NewHttpClientMock("sometoken")
			mockClient.MockAddRules = tt.mockRequest

			instance := NewRules(mockClient)
			result, err := instance.Create(tt.body, false)

			if err != nil {
				t.Errorf("got err %v", err)
			}

			if result.Data[0].Id != tt.result.Data[0].Id {
				t.Errorf("got %s, want %s", result.Data[0].Id, tt.result.Data[0].Id)
			}

			if result.Data[0].Tag != tt.result.Data[0].Tag {
				t.Errorf("got %s, want %s", result.Data[0].Tag, tt.result.Data[0].Tag)
			}

			if result.Data[0].Value != tt.result.Data[0].Value {
				t.Errorf("got %s, want %s", result.Data[0].Value, tt.result.Data[0].Value)
			}

			if result.Meta.Summary.Created != tt.result.Meta.Summary.Created {
				t.Errorf("got %d, want %d", result.Meta.Summary.Created, tt.result.Meta.Summary.Created)
			}

			if result.Meta.Summary.NotCreated != tt.result.Meta.Summary.NotCreated {
				t.Errorf("got %d, want %d", result.Meta.Summary.NotCreated, tt.result.Meta.Summary.NotCreated)
			}

			if result.Meta.Sent != tt.result.Meta.Sent {
				t.Errorf("got %s, want %s", result.Meta.Sent, tt.result.Meta.Sent)
			}
		})
	}
}


func TestDelete(t *testing.T) {

	var tests = []struct {
		body        DeleteRulesRequest
		mockRequest func(queryParams string, body string) (*http.Response, error)
		result      *TwitterRuleResponse
	}{
		{
			NewDeleteRulesRequest(123),
			func(queryParams string, bodyRequest string) (*http.Response, error) {
				json := `{
					"data": [{
						"Value": "cat has:images", 
						"Tag":"cat tweets with images", 
						"id": "123"
					}],
					"meta": {
						"sent": "today",
						"summary": {
							"created": 0,
							"not_created": 1
						}
					},
					"errors": []
				}`

				body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
				}, nil

			},
			&TwitterRuleResponse{
				Data: []DataRule{
					{
						Value: "cat has:images",
						Tag:   "cat tweets with images",
						Id:    "123",
					},
				},
				Meta: MetaRule{
					Sent: "today",
					Summary: MetaSummary{
						Created:    0,
						NotCreated: 1,
					},
				},
				Errors: nil,
			},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("TestCreate (%d)", i)

		t.Run(testName, func(t *testing.T) {
			mockClient := httpclient.NewHttpClientMock("sometoken")
			mockClient.MockAddRules = tt.mockRequest

			instance := NewRules(mockClient)
			result, err := instance.Delete(tt.body, false)

			if err != nil {
				t.Errorf("got err %v", err)
			}

			if result.Data[0].Id != tt.result.Data[0].Id {
				t.Errorf("got %s, want %s", result.Data[0].Id, tt.result.Data[0].Id)
			}

			if result.Data[0].Tag != tt.result.Data[0].Tag {
				t.Errorf("got %s, want %s", result.Data[0].Tag, tt.result.Data[0].Tag)
			}

			if result.Data[0].Value != tt.result.Data[0].Value {
				t.Errorf("got %s, want %s", result.Data[0].Value, tt.result.Data[0].Value)
			}

			if result.Meta.Summary.Created != tt.result.Meta.Summary.Created {
				t.Errorf("got %d, want %d", result.Meta.Summary.Created, tt.result.Meta.Summary.Created)
			}

			if result.Meta.Summary.NotCreated != tt.result.Meta.Summary.NotCreated {
				t.Errorf("got %d, want %d", result.Meta.Summary.NotCreated, tt.result.Meta.Summary.NotCreated)
			}

			if result.Meta.Sent != tt.result.Meta.Sent {
				t.Errorf("got %s, want %s", result.Meta.Sent, tt.result.Meta.Sent)
			}
		})
	}
}


func TestGetRules(t *testing.T) {
	var tests = []struct {
		mockRequest func() (*http.Response, error)
		result      *TwitterRuleResponse
	}{
		{
			func() (*http.Response, error) {
				json := `{
					"data": [{
						"Value": "cat has:images", 
						"Tag":"cat tweets with images", 
						"id": "123456"
					}],
					"meta": {
						"sent": "today",
						"summary": {
							"created": 0,
							"not_created": 1
						}
					}
				}`

				body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
				}, nil

			},
			&TwitterRuleResponse{
				Data: []DataRule{
					{
						Value: "cat has:images",
						Tag:   "cat tweets with images",
						Id:    "123456",
					},
				},
				Meta: MetaRule{
					Sent: "today",
					Summary: MetaSummary{
						Created:    0,
						NotCreated: 1,
					},
				},
			},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("TestGetRules (%d)", i)

		t.Run(testName, func(t *testing.T) {
			mockClient := httpclient.NewHttpClientMock("sometoken")
			mockClient.MockGetRules = tt.mockRequest

			instance := NewRules(mockClient)
			result, err := instance.Get()

			if err != nil {
				t.Errorf("got err %v", err)
			}

			if result.Data[0].Id != tt.result.Data[0].Id {
				t.Errorf("got %s, want %s", result.Data[0].Id, tt.result.Data[0].Id)
			}

			if result.Data[0].Tag != tt.result.Data[0].Tag {
				t.Errorf("got %s, want %s", result.Data[0].Tag, tt.result.Data[0].Tag)
			}

			if result.Data[0].Value != tt.result.Data[0].Value {
				t.Errorf("got %s, want %s", result.Data[0].Value, tt.result.Data[0].Value)
			}

			if result.Meta.Summary.Created != tt.result.Meta.Summary.Created {
				t.Errorf("got %d, want %d", result.Meta.Summary.Created, tt.result.Meta.Summary.Created)
			}

			if result.Meta.Summary.NotCreated != tt.result.Meta.Summary.NotCreated {
				t.Errorf("got %d, want %d", result.Meta.Summary.NotCreated, tt.result.Meta.Summary.NotCreated)
			}

			if result.Meta.Sent != tt.result.Meta.Sent {
				t.Errorf("got %s, want %s", result.Meta.Sent, tt.result.Meta.Sent)
			}
		})
	}
}
