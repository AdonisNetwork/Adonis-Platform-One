package services

import "context"

type ModelGateway interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
