package services

import (
	"context"
	"fmt"
)

type AgentService struct {
	ModelGateway ModelGateway
}

type AgentRequest struct {
	Agent  string                 `json:"agent"`
	Input  string                 `json:"input"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

type AgentResponse struct {
	Markdown string                 `json:"markdown"`
	JSON     map[string]interface{} `json:"json"`
}

// Constructor
func NewAgentService(m ModelGateway) *AgentService {
	return &AgentService{ModelGateway: m}
}

// Main entry
func (s *AgentService) Execute(ctx context.Context, req AgentRequest) (*AgentResponse, error) {
	switch req.Agent {
	case "analysis-agent":
		return s.runAnalysisAgent(ctx, req)
	case "research-agent":
		return s.runResearchAgent(ctx, req)
	case "writer-agent":
		return s.runWriterAgent(ctx, req)
	default:
		return nil, fmt.Errorf("unknown agent: %s", req.Agent)
	}
}
