package bedrock

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockClient interface {
	Converse(
		ctx context.Context,
		params *bedrockruntime.ConverseInput,
		optFns ...func(*bedrockruntime.Options),
	) (*bedrockruntime.ConverseOutput, error)
	ConverseStream(
		ctx context.Context,
		params *bedrockruntime.ConverseStreamInput,
		optFns ...func(*bedrockruntime.Options),
	) (*bedrockruntime.ConverseStreamOutput, error)
}

type Controller struct {
	client BedrockClient
}

func NewController() (Controller, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	return Controller{
		client: client,
	}, nil
}

func (c Controller) InvokeBedrock(
	ctx context.Context,
	bedrockReq bedrockruntime.ConverseInput,
) (*bedrockruntime.ConverseOutput, error) {
	output, err := c.client.Converse(ctx, &bedrockReq)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to invoke bedrock", err)
	}
	return output, nil
}

func (c Controller) StreamBedrock(
	ctx context.Context,
	bedrockReq bedrockruntime.ConverseStreamInput,
) (*bedrockruntime.ConverseStreamOutput, error) {
	output, err := c.client.ConverseStream(ctx, &bedrockReq)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to invoke bedrock", err)
	}
	return output, nil
}
