package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DefangLabs/bedrock-sidecar/convert"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockClientInterface interface {
	Converse(
		ctx context.Context,
		params *bedrockruntime.ConverseInput,
		optFns ...func(*bedrockruntime.Options),
	) (*bedrockruntime.ConverseOutput, error)
}

func invokeBedrock(
	ctx context.Context,
	bedrockReq bedrockruntime.ConverseInput,
) (*bedrockruntime.ConverseOutput, error) {
	mapped, ok := modelNameMap[*bedrockReq.ModelId]
	if ok {
		bedrockReq.ModelId = aws.String(mapped)
	}
	output, err := bedrockClient.Converse(ctx, &bedrockReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Bedrock API: %w", err)
	}

	return output, nil
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var openAIReq convert.OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bedrockReq := convert.ToBedrockRequest(openAIReq)
	bedrockResp, err := invokeBedrock(r.Context(), bedrockReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	openAIResp := convert.ToOpenAIResponse(bedrockResp, openAIReq.Model)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openAIResp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
