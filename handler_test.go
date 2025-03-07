package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DefangLabs/bedrock-sidecar/convert"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type mockBedrockClient struct {
	response *bedrockruntime.ConverseOutput
	err      error
}

func (m *mockBedrockClient) Converse(
	context.Context,
	*bedrockruntime.ConverseInput,
	...func(*bedrockruntime.Options),
) (*bedrockruntime.ConverseOutput, error) {
	return m.response, m.err
}

func TestHandleChatCompletions(t *testing.T) {
	// Save the original client and restore it after tests
	originalClient := bedrockClient
	defer func() {
		bedrockClient = originalClient
	}()

	tests := []struct {
		name         string
		method       string
		requestBody  interface{}
		mockResponse *bedrockruntime.ConverseOutput
		mockError    error
		expectedCode int
		validateResp func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful request",
			method: http.MethodPost,
			requestBody: convert.OpenAIRequest{
				Model: "gpt-3.5-turbo",
				Messages: []convert.OpenAIMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			mockResponse: &bedrockruntime.ConverseOutput{
				Output: &types.ConverseOutputMemberMessage{
					Value: types.Message{
						Content: []types.ContentBlock{
							&types.ContentBlockMemberText{
								Value: "Hi there!",
							},
						},
					},
				},
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				var resp convert.OpenAIResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Model != "gpt-3.5-turbo" {
					t.Errorf("Expected model 'gpt-3.5-turbo', got %q", resp.Model)
				}
				if len(resp.Choices) != 1 {
					t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
				}
				if resp.Choices[0].Message.Content != "Hi there!" {
					t.Errorf("Expected content 'Hi there!', got %q", resp.Choices[0].Message.Content)
				}
			},
		},
		{
			name:         "invalid method",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "invalid request body",
			method:       http.MethodPost,
			requestBody:  "invalid json",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "bedrock error",
			method: http.MethodPost,
			requestBody: convert.OpenAIRequest{
				Model: "gpt-3.5-turbo",
				Messages: []convert.OpenAIMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			mockError:    &types.ValidationException{Message: aws.String("Invalid request")},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bedrockClient = &mockBedrockClient{
				response: tt.mockResponse,
				err:      tt.mockError,
			}

			var body bytes.Buffer
			if tt.requestBody != nil {
				if err := json.NewEncoder(&body).Encode(tt.requestBody); err != nil {
					t.Fatalf("Failed to encode request body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/v1/chat/completions", &body)
			w := httptest.NewRecorder()

			handleChatCompletions(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, w)
			}
		})
	}
}
