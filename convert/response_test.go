package convert

import (
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/stretchr/testify/assert"
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
