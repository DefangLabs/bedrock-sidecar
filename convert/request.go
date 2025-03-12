package convert

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type OpenAIRequest struct {
	Model          string                 `json:"model"`
	N              int                    `json:"n"`
	MaxTokens      int                    `json:"max_tokens,omitempty"`
	ResponseFormat string                 `json:"response_format,omitempty"`
	Messages       []OpenAIMessage        `json:"messages"`
	Seed           int                    `json:"seed,omitempty"`
	Stop           []string               `json:"stop,omitempty"`
	Stream         bool                   `json:"stream,omitempty"`
	Temperature    *float64               `json:"temperature,omitempty"`
	TopP           *float64               `json:"top_p,omitempty"`
	Extra          map[string]interface{} `json:"-"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func ToBedrockRequest(openAIReq OpenAIRequest) bedrockruntime.ConverseInput {
	systemMessages := make([]string, 0, 1)
	messages := make([]types.Message, 0, len(openAIReq.Messages))

	for _, msg := range openAIReq.Messages {
		if msg.Role == "system" {
			systemMessages = append(systemMessages, msg.Content)
		} else {
			messages = append(messages, types.Message{
				Role: types.ConversationRole(msg.Role),
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: msg.Content,
					},
				},
			})
		}
	}

	var temperature *float32
	var maxTokens *int32
	var topP *float32
	if openAIReq.Temperature != nil {
		temperature = aws.Float32(float32(*openAIReq.Temperature))
	}

	if openAIReq.MaxTokens != 0 {
		maxTokens = aws.Int32(int32(openAIReq.MaxTokens)) //nolint:gosec
	}

	if openAIReq.TopP != nil {
		topP = aws.Float32(float32(*openAIReq.TopP))
	}

	return bedrockruntime.ConverseInput{
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens:     maxTokens,
			StopSequences: openAIReq.Stop,
			Temperature:   temperature,
			TopP:          topP,
		},
		Messages: messages,
		ModelId:  aws.String(openAIReq.Model),
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{
				Value: strings.Join(systemMessages, "\n"),
			},
		},
	}
}
