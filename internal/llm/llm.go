package llm

import (
	"context"

	"github.com/dipmed/backend-chat/internal/sessions"
)

type ChatRequest struct {
	Messages []sessions.Message

	// Context injected by RAG before calling the provider.
	RAGContext string
}

type Provider interface {
	ChatStream(ctx context.Context, req *ChatRequest, onChunk func(chunk string) error) error
}
