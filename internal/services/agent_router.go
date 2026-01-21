package services

import (
	"strings"
	"fmt"
)

type Domain string
type AgentName string

const (
	DomainAI       Domain = "ai"
	DomainHealth   Domain = "health"
	DomainFinance  Domain = "finance"
	DomainWeb      Domain = "web"
	DomainSecurity Domain = "security"
)

const (
	AgentAnalysis    AgentName = "analysis-agent"
	AgentResearch    AgentName = "research-agent"
	AgentWriter      AgentName = "writer-agent"
	AgentHealth      AgentName = "health-agent"
	AgentFX          AgentName = "fx-agent"
	AgentSettlement  AgentName = "settlement-agent"
	AgentCompliance  AgentName = "compliance-agent"
	AgentSEO         AgentName = "seo-agent"
	AgentSecurity    AgentName = "security-agent"
)

// Result of routing a task_type
type AgentRoute struct {
	Domain   Domain    `json:"domain"`
	Action   string    `json:"action"`
	Agent    AgentName `json:"agent"`
}

// Parse "<domain>.<action>" into parts
func parseTaskType(taskType string) (Domain, string, error) {
	parts := strings.SplitN(taskType, ".", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid task_type: %s", taskType)
	}
	return Domain(parts[0]), parts[1], nil
}

// Core router: maps domain + action → agent
func ResolveAgent(taskType string) (*AgentRoute, error) {
	domain, action, err := parseTaskType(taskType)
	if err != nil {
		return nil, err
	}

	switch domain {
	case DomainAI:
		return routeAI(domain, action), nil
	case DomainHealth:
		return routeHealth(domain, action), nil
	case DomainFinance:
		return routeFinance(domain, action), nil
	case DomainWeb:
		return routeWeb(domain, action), nil
	case DomainSecurity:
		return routeSecurity(domain, action), nil
	default:
		// fallback: treat as generic AI analysis
		return &AgentRoute{
			Domain: domain,
			Action: action,
			Agent:  AgentAnalysis,
		}, nil
	}
}

// Domain-specific routing

func routeAI(domain Domain, action string) *AgentRoute {
	switch action {
	case "analysis":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentAnalysis}
	case "research":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentResearch}
	case "writer":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentWriter}
	default:
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentAnalysis}
	}
}

func routeHealth(domain Domain, action string) *AgentRoute {
	switch action {
	case "screening":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentHealth}
	case "report":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentWriter}
	default:
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentHealth}
	}
}

func routeFinance(domain Domain, action string) *AgentRoute {
	switch action {
	case "fx_pricing":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentFX}
	case "settlement":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentSettlement}
	case "compliance_check":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentCompliance}
	default:
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentFX}
	}
}

func routeWeb(domain Domain, action string) *AgentRoute {
	switch action {
	case "seo_audit":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentSEO}
	case "content_brief":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentWriter}
	default:
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentSEO}
	}
}

func routeSecurity(domain Domain, action string) *AgentRoute {
	switch action {
	case "risk_review", "config_audit":
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentSecurity}
	default:
		return &AgentRoute{Domain: domain, Action: action, Agent: AgentSecurity}
	}
}
