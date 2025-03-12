package convert

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/stretchr/testify/assert"
)

func TestToBedrockRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    OpenAIRequest
		expected bedrockruntime.ConverseInput
	}{
		{
			name: "basic conversion",
			input: OpenAIRequest{
				Model:     "anthropic.claude-v2",
				MaxTokens: 1000,
				Messages: []OpenAIMessage{
					{Role: "system", Content: "You are a helpful assistant"},
					{Role: "user", Content: "Hello"},
				},
				Temperature: aws.Float64(0.7),
				TopP:        aws.Float64(0.9),
				Stop:        []string{"\n", "Human:"},
			},
			expected: bedrockruntime.ConverseInput{
				ModelId: aws.String("anthropic.claude-v2"),
				InferenceConfig: &types.InferenceConfiguration{
					MaxTokens:     aws.Int32(1000),
					Temperature:   aws.Float32(0.7),
					TopP:          aws.Float32(0.9),
					StopSequences: []string{"\n", "Human:"},
				},
				Messages: []types.Message{
					{
						Role: types.ConversationRole("system"),
						Content: []types.ContentBlock{
							&types.ContentBlockMemberText{
								Value: "You are a helpful assistant",
							},
						},
					},
					{
						Role: types.ConversationRole("user"),
						Content: []types.ContentBlock{
							&types.ContentBlockMemberText{
								Value: "Hello",
							},
						},
					},
				},
			},
		},
		{
			name: "empty optional fields",
			input: OpenAIRequest{
				Model: "anthropic.claude-v2",
				Messages: []OpenAIMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			expected: bedrockruntime.ConverseInput{
				ModelId: aws.String("anthropic.claude-v2"),
				InferenceConfig: &types.InferenceConfiguration{
					StopSequences: nil,
				},
				Messages: []types.Message{
					{
						Role: types.ConversationRole("user"),
						Content: []types.ContentBlock{
							&types.ContentBlockMemberText{
								Value: "Hello",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBedrockRequest(tt.input)

			assert.Equal(t, tt.expected.ModelId, result.ModelId)

			assert.Equal(t, tt.expected.InferenceConfig.MaxTokens, result.InferenceConfig.MaxTokens)
			assert.Equal(t, tt.expected.InferenceConfig.Temperature, result.InferenceConfig.Temperature)
			assert.Equal(t, tt.expected.InferenceConfig.TopP, result.InferenceConfig.TopP)
			assert.Equal(t, tt.expected.InferenceConfig.StopSequences, result.InferenceConfig.StopSequences)

			assert.Equal(t, len(tt.expected.Messages), len(result.Messages))
			for i := range tt.expected.Messages {
				assert.Equal(t, tt.expected.Messages[i].Role, result.Messages[i].Role)
				assert.Equal(t, tt.expected.Messages[i].Content[0].(*types.ContentBlockMemberText).Value,
					result.Messages[i].Content[0].(*types.ContentBlockMemberText).Value)
			}
		})
	}
}
