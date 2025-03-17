package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DefangLabs/bedrock-sidecar/bedrock"
	"github.com/DefangLabs/bedrock-sidecar/convert"
)

type Handler struct {
	Converser bedrock.BedrockConverser
	ModelMap  bedrock.ModelMap
}

func (h Handler) HandleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var openAIReq convert.OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Debug("Received", "request", openAIReq)

	if openAIReq.Stream {
		h.handleStreamedChatCompletion(ctx, w, openAIReq)
	} else {
		h.handleBufferedChatCompletion(ctx, w, openAIReq)
	}
}

func (h Handler) handleStreamedChatCompletion(
	ctx context.Context,
	w http.ResponseWriter,
	openAIReq convert.OpenAIRequest,
) {
	bedrockReq := convert.ToBedrockStreamRequest(h.ModelMap, openAIReq)
	slog.Debug("Converted", "request", bedrockReq)
	bedrockResp, err := h.Converser.ConverseStream(ctx, &bedrockReq)
	if err != nil {
		slog.Error("Failed to invoke Bedrock ConverseStream", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Debug("Received", "response", bedrockResp)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for event := range bedrockResp.GetStream().Events() {
		slog.Debug("Received", "chunk", event)
		openAIChunk := convert.ToOpenAIResponseChunk(event, openAIReq.Model)
		slog.Debug("Converted", "chunk", openAIChunk)
		data, err := json.Marshal(openAIChunk)
		if err != nil {
			slog.Error("Failed to encode response", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		message := []byte(fmt.Sprintf("data: %s\n\n", data))
		if _, err := w.Write(message); err != nil {
			slog.Error("Failed to write data", "error", err)
			http.Error(w, "Failed to write data", http.StatusInternalServerError)
			return
		}
		flusher.Flush()
	}
}

func (h Handler) handleBufferedChatCompletion(
	ctx context.Context,
	w http.ResponseWriter,
	openAIReq convert.OpenAIRequest,
) {
	bedrockReq := convert.ToBedrockRequest(h.ModelMap, openAIReq)
	slog.Debug("Converted", "request", bedrockReq)
	bedrockResp, err := h.Converser.Converse(ctx, &bedrockReq)
	if err != nil {
		slog.Error("Failed to invoke Bedrock Converse", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Debug("Received", "response", bedrockResp)
	openAIResp := convert.ToOpenAIResponse(bedrockResp, openAIReq.Model)
	slog.Debug("Converted", "response", bedrockResp)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openAIResp); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
