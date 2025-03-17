# Bedrock Sidecar

A proxy server that converts OpenAI API requests to AWS Bedrock requests, allowing you to use OpenAI-compatible clients with AWS Bedrock.
The Defang Bedrock Sidecar is deployed alongside Defang services, and presents an OpenAI-compatible interface for AWS Bedrock.

## Usage

### Enable model access in Bedrock

1. Navigate to the [Bedrock Model Access dashboard](https://console.aws.amazon.com/bedrock/home#/modelaccess)
1. Click the "Modify model access" button
1. Select the models you would like to use with Bedrock
1. Click the "Next" button at the bottom of the page
1. Click the "Submit" button

### Update your compose file

```yaml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    # Add this tag to set up integration with bedrock when deploying to aws
    x-defang-llm: true
```

## Features

- Intercepts OpenAI API requests
- Converts OpenAI chat completion requests to AWS Bedrock format
- Converts AWS Bedrock responses back to OpenAI format
- Supports basic chat completion functionality

## Limitations

- Currently only supports basic chat completion functionality
- Token counting is not implemented

## Contributing

See our [Contributing](./CONTRIBUTING.md) document.
