package main

import (
	"log"
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
	bedrockController, err := bedrock.NewController()
	if err != nil {
		log.Fatalf("unable to create bedrock controller: %v", err)
	}

	modelMap, err := bedrock.NewModelMap()
	if err != nil {
		log.Fatalf("unable to create model map: %v", err)
	}

	handler := handler.Handler{
		Converser: bedrockController,
		ModelMap:  modelMap,
	}

	http.HandleFunc("/v1/chat/completions", handler.HandleChatCompletions)
	http.HandleFunc("/api/chat", handler.HandleChatCompletions)

	log.Printf("Listening on port %s", port)

	srv := &http.Server{
		Addr:              ":" + port,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
