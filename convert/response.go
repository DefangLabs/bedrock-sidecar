package convert

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type TimeProvider func() time.Time

var timeProvider TimeProvider = time.Now

func SetTimeProvider(provider TimeProvider) TimeProvider {
	old := timeProvider
	timeProvider = provider
	return old
}

func ToOpenAIResponse(bedrockOutput *bedrockruntime.ConverseOutput, model string) OpenAIResponse {
	now := timeProvider()

	var message types.ConverseOutputMemberMessage

	if msg, ok := bedrockOutput.Output.(*types.ConverseOutputMemberMessage); ok {
		message = *msg
	} else {
		return OpenAIResponse{
			ID:      generateID(),
			Object:  "chat.completion",
			Created: now.Unix(),
			Model:   model,
			Choices: []Choice{
				{
					Index: 0,
					Message: OpenAIMessage{
						Role:    "system",
						Content: "Error: invalid message type",
					},
					FinishReason: "error",
				},
			},
			Usage: Usage{
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			},
		}
	}
	var content []string
	for _, cntnt := range message.Value.Content {
		if textContent, ok := cntnt.(*types.ContentBlockMemberText); ok {
			content = append(content, textContent.Value)
			break
		}
	}

	return OpenAIResponse{
		ID:      generateID(),
		Object:  "chat.completion",
		Created: now.Unix(),
		Model:   model,
		Choices: []Choice{
			{
				Index: 0,
				Message: OpenAIMessage{
					Role:    "assistant",
					Content: strings.Join(content, "\n\n"),
				},
				FinishReason: string(bedrockOutput.StopReason),
			},
		},
		Usage: Usage{
			PromptTokens:     0, // TODO: Implement token counting
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}
}

func generateID() string {
	return "chatcmpl-" + strconv.FormatInt(timeProvider().UnixNano(), 10)
}
