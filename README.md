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

## Prerequisites

- Go 1.22 or later
- AWS credentials configured (either through environment variables, AWS CLI, or IAM role)
- AWS Bedrock access configured in your AWS account

## Setup

1. Clone the repository
1. Build the project
   ```bash
   make build
   ```
1. Run the tests
   ```bash
   make test
   ```

## Configuration

The following environment variables can be used to configure the proxy:

- `PORT`: The port number to run the server on (default: 8080)
- Standard AWS configuration environment variables (AWS_REGION, AWS_ACCESS_KEY_ID, etc.)

## Running the Server

```bash
make run
```

The server will start on port 8080 (or the port specified in the PORT environment variable).

## Usage

The proxy implements the OpenAI chat completions endpoint. You can use it as a drop-in replacement for the OpenAI API by pointing your OpenAI-compatible client to:

```
http://localhost:8080/v1/chat/completions
```

Example curl request:

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "meta.llama3-8b-instruct-v1:0",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ],
    "max_tokens": 2048
  }'
```

## Limitations

- Currently only supports basic chat completion functionality
- Streaming is not yet implemented
- Token counting is not implemented

## Make Docker images

```bash
make images
```

## Push Docker images

```bash
make push-images
```
