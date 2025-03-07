package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

var bedrockClient BedrockClientInterface

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	bedrockClient = bedrockruntime.NewFromConfig(cfg)

	http.HandleFunc("/v1/chat/completions", handleChatCompletions)

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
