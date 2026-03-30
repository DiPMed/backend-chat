# DiPMed Backend Chat

LLM-powered chat backend with streaming and session support.

## Prerequisites

- Go 1.22+
- [Gemini API key](https://aistudio.google.com/apikey)

## Build & Run

```bash
export LLM_API_KEY="your-gemini-api-key"
make build
make run
```

## Example usage in terminal

Start a conversation:

```bash
curl -N -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello"}'
```

Continue a conversation (use `session_id` from the first response chunk):

```bash
curl -N -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"session_id": "<id>", "message": "Tell me more"}'
```
