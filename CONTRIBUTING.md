
# Contributing

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
- `MODEL_NAME_MAP`: A json object string which maps an openai model name to a bedrock model name. For example: `MODEL_NAME_MAP='{"gpt-4o": "anthropic.claude-3-5-sonnet-20241022-v2:0"}'`
- Standard AWS configuration environment variables (AWS_REGION, AWS_ACCESS_KEY_ID, etc.)

## Running the server

```bash
make run
```

The server will start on port 8080 (or the port specified in the PORT environment variable).

## Example request

The proxy implements the OpenAI chat completions endpoint. You can use it as a drop-in replacement for the OpenAI API by pointing your OpenAI-compatible client to:

```
http://localhost:8080/v1/chat/completions
```

Example curl request:

```bash
$ AWS_PROFILE=foo AWS_REGION=us-west-2 MODEL_NAME_MAP='{"gpt-4o": "anthropic.claude-3-5-sonnet-20241022-v2:0"}' make run
$ curl -s -X POST http://localhost:8080/v1/chat/completions -H "Content-Type: application/json" -d '{
    "model": "gpt-4o",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ],
    "max_tokens": 2048
  }'
{
  "id": "chatcmpl-1742229889235795000",
  "object": "chat.completion",
  "created": 1742229889,
  "model": "gpt-4o",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hi! I'm doing well, thanks for asking. I'm ready to help in whatever way I can. How are you today?"
      },
      "finish_reason": "end_turn"
    }
  ],
  "usage": {
    "prompt_tokens": 0,
    "completion_tokens": 0,
    "total_tokens": 0
  }
}
```

## Docker images

Defang publishes a docker image to [`defangio/bedrock-sidecar`](https://hub.docker.com/r/defangio/bedrock-sidecar).

### Making Docker images

```bash
make images
```

### Pushing Docker images

```bash
make push-images
```
