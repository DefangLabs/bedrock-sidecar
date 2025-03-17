package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/DefangLabs/bedrock-sidecar/bedrock"
	"github.com/DefangLabs/bedrock-sidecar/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	debug := os.Getenv("DEBUG")
	if debug != "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}

	bedrockController, err := bedrock.NewController()
	if err != nil {
		slog.Error("Failed to create bedrock.Controller", "error", err)
		os.Exit(1)
	}

	modelMap, err := bedrock.NewModelMap()
	if err != nil {
		slog.Error("Failed to create bedrock.ModelMap", "error", err)
		os.Exit(1)
	}

	handler := handler.Handler{
		Converser: bedrockController,
		ModelMap:  modelMap,
	}

	http.HandleFunc("/v1/chat/completions", handler.HandleChatCompletions)
	http.HandleFunc("/api/chat", handler.HandleChatCompletions)

	slog.Info("Listening", "port", port)

	srv := &http.Server{
		Addr:              ":" + port,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}
}
