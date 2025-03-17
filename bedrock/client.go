package bedrock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockConverser interface {
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

type Client struct {
	client BedrockConverser
}

func NewController() (Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return Client{}, fmt.Errorf("failed to load SDK config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	return Client{
		client: client,
	}, nil
}

func (c Client) Converse(
	ctx context.Context,
	bedrockReq *bedrockruntime.ConverseInput,
	optFns ...func(*bedrockruntime.Options),
) (*bedrockruntime.ConverseOutput, error) {
	output, err := c.client.Converse(ctx, bedrockReq, optFns...)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to invoke bedrock", err)
	}
	return output, nil
}

func (c Client) ConverseStream(
	ctx context.Context,
	bedrockReq *bedrockruntime.ConverseStreamInput,
	optFns ...func(*bedrockruntime.Options),
) (*bedrockruntime.ConverseStreamOutput, error) {
	output, err := c.client.ConverseStream(ctx, bedrockReq, optFns...)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to invoke bedrock", err)
	}
	return output, nil
}
