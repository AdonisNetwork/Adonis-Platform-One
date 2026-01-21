package services

import (
	"context"
	"fmt"
)

type OpenAIModelGateway struct {
	APIKey string
}

func NewOpenAIModelGateway(apiKey string) *OpenAIModelGateway {
	return &OpenAIModelGateway{APIKey: apiKey}
}

func (g *OpenAIModelGateway) Generate(ctx context.Context, prompt string) (string, error) {
	// واقعی بعداً
	return fmt.Sprintf("LLM Output for prompt: %s", prompt), nil
}
