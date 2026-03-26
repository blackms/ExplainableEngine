package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/blackms/ExplainableEngine/internal/models"
)

// ClaudeService implements the Service interface using the Anthropic Claude API.
type ClaudeService struct {
	client anthropic.Client
	model  anthropic.Model
}

// NewClaudeService creates a new ClaudeService with the given API key.
func NewClaudeService(apiKey string) (*ClaudeService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &ClaudeService{
		client: client,
		model:  anthropic.ModelClaudeSonnet4_20250514,
	}, nil
}

// IsAvailable returns true if the Claude client is initialized.
func (s *ClaudeService) IsAvailable() bool { return true }

// GenerateNarrative asks Claude to produce a human-readable narrative from the
// given explanation data.
func (s *ClaudeService) GenerateNarrative(ctx context.Context, explanation *models.ExplainResponse, level, lang string) (string, error) {
	systemPrompt := buildNarrativeSystemPrompt(level, lang)
	userPrompt := BuildExplanationContext(explanation)

	resp, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     s.model,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude API error: %w", err)
	}

	return extractText(resp)
}

// AnswerQuestion provides a conversational Q&A interface grounded on the
// explanation data.
func (s *ClaudeService) AnswerQuestion(ctx context.Context, explanation *models.ExplainResponse, question string, history []Message) (string, error) {
	systemPrompt := "You are an AI assistant that answers questions about numerical explanations. " +
		"You have access to the full breakdown, dependency graph, confidence scores, and driver analysis. " +
		"Answer based ONLY on the data provided. If you cannot answer from the data, say so."

	contextText := BuildExplanationContext(explanation)

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("Here is the explanation data:\n\n" + contextText)),
		anthropic.NewAssistantMessage(anthropic.NewTextBlock("I have reviewed the explanation data. What would you like to know?")),
	}

	// Append conversation history.
	for _, msg := range history {
		if msg.Role == "user" {
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
		} else {
			messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
		}
	}

	// Append current question.
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(question)))

	resp, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     s.model,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("claude API error: %w", err)
	}

	return extractText(resp)
}

// GenerateSummary produces an audience-specific executive summary as structured
// JSON from the explanation data.
func (s *ClaudeService) GenerateSummary(ctx context.Context, explanation *models.ExplainResponse, audience, lang string) (*SummaryResult, error) {
	systemPrompt := buildSummarySystemPrompt(audience, lang)
	userPrompt := BuildExplanationContext(explanation)

	resp, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     s.model,
		MaxTokens: 2048,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	text, err := extractText(resp)
	if err != nil {
		return nil, err
	}

	// Strip markdown code fences if present.
	text = stripCodeFences(text)

	var result SummaryResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse summary JSON: %w (raw: %s)", err, text)
	}
	result.Audience = audience
	result.Language = lang
	return &result, nil
}

// extractText pulls the first text block from a Claude API response.
func extractText(resp *anthropic.Message) (string, error) {
	for _, block := range resp.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}
	return "", fmt.Errorf("no text in response")
}

// stripCodeFences removes leading/trailing markdown code fences (```json ... ```).
func stripCodeFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		// Remove first line (```json or ```)
		if idx := strings.Index(s, "\n"); idx >= 0 {
			s = s[idx+1:]
		}
	}
	if strings.HasSuffix(s, "```") {
		s = s[:len(s)-3]
	}
	return strings.TrimSpace(s)
}

// buildNarrativeSystemPrompt returns the system prompt for narrative generation
// tailored to the requested detail level and language.
func buildNarrativeSystemPrompt(level, lang string) string {
	langInstr := "Respond in English."
	if lang == "it" {
		langInstr = "Rispondi in italiano."
	}

	switch level {
	case "executive":
		return fmt.Sprintf("You are a financial analyst writing an executive briefing. Be concise, authoritative, and focus on business implications. %s", langInstr)
	case "advanced":
		return fmt.Sprintf("You are a quantitative analyst explaining a model's output. Include technical details about component contributions, confidence levels, and data quality. %s", langInstr)
	default: // basic
		return fmt.Sprintf("You are explaining a numerical result to a non-technical stakeholder. Be clear, simple, and reassuring. One paragraph maximum. %s", langInstr)
	}
}

// buildSummarySystemPrompt returns the system prompt for structured summary
// generation based on the target audience.
func buildSummarySystemPrompt(audience, lang string) string {
	langInstr := "Respond in English."
	if lang == "it" {
		langInstr = "Rispondi in italiano."
	}

	base := `You must respond ONLY with a JSON object matching this schema:
{"title":"string","summary":"string","key_findings":["string"],"risks":["string"],"recommendations":["string"]}
Do not include any text outside the JSON object.`

	switch audience {
	case "board":
		return fmt.Sprintf("You are preparing an executive board briefing. Focus on high-level business impact. Avoid jargon. %s\n\n%s", langInstr, base)
	case "technical":
		return fmt.Sprintf("You are preparing a technical deep-dive report. Include full technical detail, formulas, and component analysis. %s\n\n%s", langInstr, base)
	default: // client
		return fmt.Sprintf("You are preparing a client-facing report. Use a balanced, trustworthy tone. Focus on what the result means for them. %s\n\n%s", langInstr, base)
	}
}

// BuildExplanationContext serializes an ExplainResponse into a clear text
// representation suitable as LLM context.
func BuildExplanationContext(e *models.ExplainResponse) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Target: %s\n", e.Target))
	sb.WriteString(fmt.Sprintf("Final Value: %.4f\n", e.FinalValue))
	sb.WriteString(fmt.Sprintf("Overall Confidence: %.1f%%\n", e.Confidence*100))
	sb.WriteString(fmt.Sprintf("Missing Data Impact: %.1f%%\n", e.MissingImpact*100))

	if len(e.Breakdown) > 0 {
		sb.WriteString("\nBreakdown:\n")
		for _, b := range e.Breakdown {
			sb.WriteString(fmt.Sprintf("  - %s: value=%.4f, weight=%.2f, contribution=%.4f (%.1f%%), confidence=%.1f%%\n",
				b.Label, b.Value, b.Weight, b.AbsoluteContribution, b.Percentage, b.Confidence*100))
		}
	}

	if len(e.TopDrivers) > 0 {
		sb.WriteString("\nTop Drivers:\n")
		for _, d := range e.TopDrivers {
			sb.WriteString(fmt.Sprintf("  #%d %s (impact: %.4f)\n", d.Rank, d.Name, d.Impact))
		}
	}

	if e.Graph != nil {
		sb.WriteString(fmt.Sprintf("\nGraph: %d nodes, %d edges\n", len(e.Graph.Nodes), len(e.Graph.Edges)))
	}

	if e.DependencyTree != nil {
		sb.WriteString(fmt.Sprintf("Dependency Tree: depth=%d, total_nodes=%d\n", e.DependencyTree.Depth, e.DependencyTree.TotalNodes))
	}

	return sb.String()
}
