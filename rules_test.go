package twitterstream

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAddRules(t *testing.T) {

	var tests = []struct {
		body        string
		mockRequest func(opts *requestOpts) (*http.Response, error)
		result      *rulesResponse
	}{
		{
			`{
				"add": [
					{"value": "cat has:images", "tag": "cat tweets with images"}
				]
			}`,
			func(opts *requestOpts) (*http.Response, error) {
				json := `{
					"data": [{
						"value": "cat has:images", 
						"tag":"cat tweets with images", 
						"id": "123456"
					}],
					"meta": {
						"sent": "today",
						"summary": {
							"created": 1,
							"not_created": 0
						}
					}
				}`

				body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
				}, nil

			},
			&rulesResponse{
				Data: []rulesResponseValue{
					{
						Value: "cat has:images",
						Tag:   "cat tweets with images",
						Id:    "123456",
					},
				},
				Meta: rulesResponseMeta{
					Sent: "today",
					Summary: addRulesResponseMetaSummary{
						Created:    1,
						NotCreated: 0,
					},
				},
			},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("TestAddRules (%d) %s", i, tt.body)

		t.Run(testName, func(t *testing.T) {
			mockClient := newHttpClientMock("sometoken")
			mockClient.MockNewHttpRequest = tt.mockRequest

			instance := newRules(mockClient)
			result, err := instance.AddRules(tt.body, false)

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
		mockRequest func(opts *requestOpts) (*http.Response, error)
		result      *rulesResponse
	}{
		{
			func(opts *requestOpts) (*http.Response, error) {
				json := `{
					"data": [{
						"value": "cat has:images", 
						"tag":"cat tweets with images", 
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
			&rulesResponse{
				Data: []rulesResponseValue{
					{
						Value: "cat has:images",
						Tag:   "cat tweets with images",
						Id:    "123456",
					},
				},
				Meta: rulesResponseMeta{
					Sent: "today",
					Summary: addRulesResponseMetaSummary{
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
			mockClient := newHttpClientMock("sometoken")
			mockClient.MockNewHttpRequest = tt.mockRequest

			instance := newRules(mockClient)
			result, err := instance.GetRules()

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