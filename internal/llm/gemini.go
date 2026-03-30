package llm

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
	model  string
}

func NewGeminiProvider(ctx context.Context, apiKey, model string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating gemini client: %w", err)
	}

	return &GeminiProvider{client: client, model: model}, nil
}

func (g *GeminiProvider) ChatStream(ctx context.Context, req *ChatRequest, onChunk func(string) error) error {
	contents := make([]*genai.Content, 0, len(req.Messages))
	for _, m := range req.Messages {
		contents = append(contents, &genai.Content{
			Role:  mapRole(m.Role),
			Parts: []*genai.Part{genai.NewPartFromText(m.Content)},
		})
	}

	var systemParts []*genai.Part
	if req.RAGContext != "" {
		systemParts = append(systemParts, genai.NewPartFromText(req.RAGContext))
	}

	last := contents[len(contents)-1]
	prev := contents[:len(contents)-1]

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{Parts: systemParts},
	}

	chat, err := g.client.Chats.Create(ctx, g.model, config, prev)
	if err != nil {
		return fmt.Errorf("creating chat session: %w", err)
	}

	lastParts := make([]genai.Part, len(last.Parts))
	for i, p := range last.Parts {
		lastParts[i] = *p
	}

	for resp, err := range chat.SendMessageStream(ctx, lastParts...) {
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}
		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				if part.Text != "" {
					if err := onChunk(part.Text); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func mapRole(role string) string {
	switch role {
	case "assistant":
		return "model"
	default:
		return role
	}
}
