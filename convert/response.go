package convert

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/openai/openai-go"
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

func ToOpenAIResponseChunk(bedrockChunk types.ConverseStreamOutput, model string) openai.ChatCompletionChunk {
	now := timeProvider()

	choice := makeOpenAIChatCompletionChunkChoice(bedrockChunk)

	return openai.ChatCompletionChunk{
		ID:      generateID(),
		Object:  "chat.completion.chunk",
		Created: now.Unix(),
		Model:   model,
		Choices: []openai.ChatCompletionChunkChoice{
			choice,
		},
	}
}

func makeOpenAIChatCompletionChunkChoice(bedrockChunk types.ConverseStreamOutput) openai.ChatCompletionChunkChoice {
	choice := openai.ChatCompletionChunkChoice{}

	switch output := bedrockChunk.(type) {
	case *types.ConverseStreamOutputMemberContentBlockStart:
		log.Println("handling of ConverseStreamOutputMemberContentBlockStart in unimplemented")
	case *types.ConverseStreamOutputMemberContentBlockStop:
		log.Println("handling of ConverseStreamOutputMemberContentBlockStop in unimplemented")
	case *types.ConverseStreamOutputMemberMetadata:
		log.Println("handling of ConverseStreamOutputMemberMetadata in unimplemented")
	case *types.ConverseStreamOutputMemberMessageStart:
		choice.Delta = openai.ChatCompletionChunkChoicesDelta{
			Role: openai.ChatCompletionChunkChoicesDeltaRole(output.Value.Role),
		}
	case *types.ConverseStreamOutputMemberMessageStop:
		choice.FinishReason = mapStopReasonToFinishReason(output.Value.StopReason)
	case *types.ConverseStreamOutputMemberContentBlockDelta:
		choice = handleContentBlockDelta(output)
	default:
		log.Println("union is nil or unknown type")
	}

	return choice
}

func mapStopReasonToFinishReason(stopReason types.StopReason) openai.ChatCompletionChunkChoicesFinishReason {
	switch stopReason {
	case types.StopReasonEndTurn, types.StopReasonStopSequence:
		return "stop"
	case types.StopReasonMaxTokens:
		return "length"
	case types.StopReasonContentFiltered, types.StopReasonGuardrailIntervened:
		return "content_filter"
	case types.StopReasonToolUse:
		return "tool_calls"
	default:
		return "stop"
	}
}

func handleContentBlockDelta(
	output *types.ConverseStreamOutputMemberContentBlockDelta,
) openai.ChatCompletionChunkChoice {
	choice := openai.ChatCompletionChunkChoice{}
	switch delta := output.Value.Delta.(type) {
	case *types.ContentBlockDeltaMemberText:
		choice.Delta = openai.ChatCompletionChunkChoicesDelta{
			Content: delta.Value,
		}
	case *types.ContentBlockDeltaMemberReasoningContent:
		log.Println("handling of ContentBlockDeltaMemberReasoningContent in unimplemented")
	case *types.ContentBlockDeltaMemberToolUse:
		log.Println("handling of ContentBlockDeltaMemberReasoningContent in unimplemented")
	}
	return choice
}

func generateID() string {
	return "chatcmpl-" + strconv.FormatInt(timeProvider().UnixNano(), 10)
}
