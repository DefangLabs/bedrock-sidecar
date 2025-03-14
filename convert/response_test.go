package convert

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToOpenAIResponse(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	oldProvider := SetTimeProvider(func() time.Time {
		return fixedTime
	})
	defer SetTimeProvider(oldProvider)

	t.Run("successful conversion", func(t *testing.T) {
		bedrockOutput := &bedrockruntime.ConverseOutput{
			Output: &types.ConverseOutputMemberMessage{
				Value: types.Message{
					Content: []types.ContentBlock{
						&types.ContentBlockMemberText{
							Value: "Test response",
						},
					},
				},
			},
			StopReason: "stop",
		}

		result := ToOpenAIResponse(bedrockOutput, "anthropic.claude-v2")

		assert.Equal(t, "chatcmpl-"+strconv.FormatInt(fixedTime.UnixNano(), 10), result.ID)
		assert.Equal(t, "chat.completion", result.Object)
		assert.Equal(t, fixedTime.Unix(), result.Created)
		assert.Equal(t, "anthropic.claude-v2", result.Model)
		assert.Len(t, result.Choices, 1)
		assert.Equal(t, 0, result.Choices[0].Index)
		assert.Equal(t, "Test response", result.Choices[0].Message.Content)
		assert.Equal(t, "assistant", result.Choices[0].Message.Role)
		assert.Equal(t, "stop", result.Choices[0].FinishReason)

		bytes, err := json.Marshal(result)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"id": "chatcmpl-1704067200000000000",
			"object": "chat.completion",
			"created": 1704067200,
			"model": "anthropic.claude-v2",
			"choices": [
				{
					"index": 0,
					"message": {

						"role": "assistant",
						"content": "Test response"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 0,
				"completion_tokens": 0,
				"total_tokens": 0
			}
		}`, string(bytes))
	})

	t.Run("error case - invalid message type", func(t *testing.T) {
		bedrockOutput := &bedrockruntime.ConverseOutput{
			Output: nil,
		}

		result := ToOpenAIResponse(bedrockOutput, "anthropic.claude-v2")

		assert.Equal(t, "chat.completion", result.Object)
		assert.Equal(t, "anthropic.claude-v2", result.Model)
		assert.Len(t, result.Choices, 1)
		assert.Equal(t, "Error: invalid message type", result.Choices[0].Message.Content)
		assert.Equal(t, "system", result.Choices[0].Message.Role)
		assert.Equal(t, "error", result.Choices[0].FinishReason)
	})
}

func TestToOpenAIResponseChunk(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	oldProvider := SetTimeProvider(func() time.Time {
		return fixedTime
	})
	defer SetTimeProvider(oldProvider)

	t.Run("successful conversion", func(t *testing.T) {
		event := &types.ConverseStreamOutputMemberContentBlockDelta{
			Value: types.ContentBlockDeltaEvent{
				ContentBlockIndex: aws.Int32(0),
				Delta: &types.ContentBlockDeltaMemberText{
					Value: "Test response",
				},
			},
		}

		result := ToOpenAIResponseChunk(event, "anthropic.claude-v2")

		assert.Equal(t, "chatcmpl-"+strconv.FormatInt(fixedTime.UnixNano(), 10), result.ID)
		assert.Equal(t, openai.ChatCompletionChunkObject("chat.completion.chunk"), result.Object)
		assert.Equal(t, fixedTime.Unix(), result.Created)
		assert.Equal(t, "anthropic.claude-v2", result.Model)
		assert.Equal(t, "Test response", result.Choices[0].Delta.Content)
		assert.Equal(t, openai.ChatCompletionChunkChoicesDeltaRole(""), result.Choices[0].Delta.Role)

		bytes, err := json.Marshal(result)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"id": "chatcmpl-1704067200000000000",
			"choices": [{
				"delta": {
					"content": "Test response",
					"function_call": { "arguments": "", "name": "" },
					"refusal": "",
					"role": "",
					"tool_calls": null
				},
				"finish_reason": "",
				"index": 0,
				"logprobs": {
					"content": null,
					"refusal": null
				}
			}],
			"created": 1704067200,
			"model": "anthropic.claude-v2",
			"object": "chat.completion.chunk",
			"service_tier": "",
			"system_fingerprint": "",
			"usage": {
				"completion_tokens": 0,
				"total_tokens": 0,
				"prompt_tokens": 0,
				"completion_tokens_details": {
					"accepted_prediction_tokens": 0,
					"audio_tokens": 0,
					"reasoning_tokens": 0,
					"rejected_prediction_tokens": 0
				},
				"prompt_tokens_details": {
					"audio_tokens": 0,
					"cached_tokens": 0
				}
			}
		}`, string(bytes))
	})
}

func TestSetTimeProvider(t *testing.T) {
	originalTime := timeProvider()
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	oldProvider := SetTimeProvider(func() time.Time {
		return fixedTime
	})

	assert.Equal(t, fixedTime, timeProvider())

	SetTimeProvider(oldProvider)

	newTime := timeProvider()
	assert.NotEqual(t, fixedTime, newTime)
	assert.True(t, newTime.After(originalTime))
}
