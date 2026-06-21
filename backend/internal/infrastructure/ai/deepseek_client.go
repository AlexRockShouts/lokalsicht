package ai

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type DeepSeekClient struct {
	client *openai.Client
}

func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.deepseek.com"
	return &DeepSeekClient{client: openai.NewClientWithConfig(config)}
}

var promptTemplate = `Du bist ein KI-Assistent für ein Schweizer KMU-Tool namens "Lokalsicht".
Du hilfst lokalen Unternehmen, auf Google-Bewertungen zu antworten.

SCHREIBE NUR DIE ANTWORT, KEINE EINLEITUNG, KEINE ERLÄUTERUNG.

Das Unternehmen: %s
Branche: %s

Die Bewertung (Sprache: %s):
Sterne: %d/5
Text: "%s"

Schreibe eine professionelle, freundliche Antwort in der Sprache %s.
Wichtige Regeln:
- Maximal 3 Sätze.
- Bedanke dich bei positiven Bewertungen.
- Bei negativen Bewertungen: entschuldige dich, bleibe sachlich, biete Lösungen an.
- Verwende NICHT "Sehr geehrte/r" — bleibe persönlich.
- Keine erfundenen Fakten über das Unternehmen.`

func (c *DeepSeekClient) GenerateReply(ctx context.Context, reviewText string, language string, businessContext string) ([]string, error) {
	// Detect language if not provided
	if language == "" {
		language = detectLanguage(reviewText)
	}

	prompt := fmt.Sprintf(promptTemplate, businessContext, "Lokales Unternehmen", language, 5, reviewText, language)

	// Generate two variants: friendly and professional
	variants := make([]string, 2)

	for i, tone := range []string{"freundlich", "professionell"} {
		fullPrompt := prompt + fmt.Sprintf("\n\nTonfall: %s", tone)

		resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: "deepseek-chat",
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Du bist ein professioneller Bewertungsmanager für lokale Unternehmen in der Schweiz."},
				{Role: openai.ChatMessageRoleUser, Content: fullPrompt},
			},
			MaxTokens:   200,
			Temperature: 0.7,
		})
		if err != nil {
			return nil, fmt.Errorf("deepseek api error: %w", err)
		}
		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no response from AI")
		}
		text := strings.TrimSpace(resp.Choices[0].Message.Content)
		variants[i] = text
	}

	return variants, nil
}

// detectLanguage uses a simple heuristic based on common words.
func detectLanguage(text string) string {
	text = strings.ToLower(text)
	deWords := []string{"und", "die", "der", "das", "ist", "sehr", "nicht", "gut", "war"}
	frWords := []string{"et", "très", "bien", "pas", "est", "pour", "dans", "avec", "une"}
	itWords := []string{"e", "molto", "bene", "non", "sono", "per", "una", "con"}

	words := strings.Fields(text)
	deCount, frCount, itCount := 0, 0, 0
	for _, w := range words {
		for _, dw := range deWords {
			if w == dw {
				deCount++
			}
		}
		for _, fw := range frWords {
			if w == fw {
				frCount++
			}
		}
		for _, iw := range itWords {
			if w == iw {
				itCount++
			}
		}
	}

	if deCount > frCount && deCount > itCount {
		return "de"
	}
	if frCount > deCount && frCount > itCount {
		return "fr"
	}
	if itCount > deCount && itCount > frCount {
		return "it"
	}
	return "de" // default
}
