package convert

import (
	"github.com/DefangLabs/bedrock-sidecar/bedrock"
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

func ToBedrockRequest(modelMap bedrock.ModelMap, openAIReq OpenAIRequest) bedrockruntime.ConverseInput {
	systemMessages, messages := partitionSystemMessages(openAIReq.Messages)

	return bedrockruntime.ConverseInput{
		InferenceConfig: makeInferenceConfig(openAIReq),
		Messages:        messages,
		ModelId:         aws.String(modelMap.BedrockModelID(openAIReq.Model)),
		System:          makeSystem(systemMessages),
	}
}

func ToBedrockStreamRequest(modelMap bedrock.ModelMap, openAIReq OpenAIRequest) bedrockruntime.ConverseStreamInput {
	systemMessages, messages := partitionSystemMessages(openAIReq.Messages)

	return bedrockruntime.ConverseStreamInput{
		InferenceConfig: makeInferenceConfig(openAIReq),
		Messages:        messages,
		ModelId:         aws.String(modelMap.BedrockModelID(openAIReq.Model)),
		System:          makeSystem(systemMessages),
	}
}

func partitionSystemMessages(openAIMessages []OpenAIMessage) ([]types.Message, []types.Message) {
	systemMessages := make([]types.Message, 0, 1)
	messages := make([]types.Message, 0, len(openAIMessages))

	for _, msg := range openAIMessages {
		bedrockMessage := types.Message{
			Role: types.ConversationRole(msg.Role),
			Content: []types.ContentBlock{
				&types.ContentBlockMemberText{
					Value: msg.Content,
				},
			},
		}

		if msg.Role == "system" {
			systemMessages = append(systemMessages, bedrockMessage)
		} else {
			messages = append(messages, bedrockMessage)
		}
	}

	return systemMessages, messages
}

func makeSystem(systemMessages []types.Message) []types.SystemContentBlock {
	system := make([]types.SystemContentBlock, 0, len(systemMessages))

	for _, msg := range systemMessages {
		system = append(system, &types.SystemContentBlockMemberText{
			Value: msg.Content[0].(*types.ContentBlockMemberText).Value,
		})
	}

	return system
}

func makeInferenceConfig(openAIReq OpenAIRequest) *types.InferenceConfiguration {
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

	return &types.InferenceConfiguration{
		MaxTokens:     maxTokens,
		StopSequences: openAIReq.Stop,
		Temperature:   temperature,
		TopP:          topP,
	}
}
