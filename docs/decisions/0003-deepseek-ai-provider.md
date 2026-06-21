---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Use DeepSeek V3 as AI provider for review responses

## Context and Problem Statement

The app needs AI-generated responses to customer reviews. Requirements: good German/French/Italian text quality, OpenAI-compatible API (to reuse `go-openai` SDK), pay-per-use pricing below CHF 10/month at MVP scale.

## Decision Drivers

- German/French/Italian language quality (primary target: Swiss market)
- OpenAI-compatible API (reuse `sashabaranov/go-openai`)
- Cost: predictably low at MVP scale (~500 generations/month)
- GDPR/data privacy considerations

## Considered Options

- DeepSeek V3 (chosen)
- Grok / xAI (rejected)
- OpenAI GPT-4o (rejected)
- Anthropic Claude Sonnet (rejected)

## Decision Outcome

**Chosen: DeepSeek V3**, accessed via `sashabaranov/go-openai` with custom `BaseURL = "https://api.deepseek.com"`.

### Consequences

- Good, because very low cost (~CHF 0.15/1M input, ~CHF 0.30/1M output)
- Good, because OpenAI-compatible API — no SDK change
- Good, because better German text quality than Grok
- Bad, because servers in China raise data privacy questions (mitigated: only public review text + business name are sent, no personal user data)
- Bad, because occasional API instability (mitigated: graceful error handling in the Go service, timeout, retry)

## Implementation Plan

- **Affected paths**: `backend/internal/infrastructure/ai/`
- **Dependencies**: `github.com/sashabaranov/go-openai`
- **Patterns to follow**: All AI calls through `ai.Client` interface
- **Patterns to avoid**: Direct OpenAI calls, hardcoded API keys outside config
- **Configuration**: `DEEPSEEK_API_KEY` env var, `AI_PROVIDER` flag for future switching

### Verification

- [ ] `go-openai` client configured with `BaseURL = "https://api.deepseek.com"`
- [ ] AI generates a German response for a German-language review
- [ ] AI generates a French response for a French-language review
- [ ] Language detection uses heuristic (word frequency), not AI inference
- [ ] Prompt is conservative — only fills template gaps, does not invent facts
- [ ] API timeout configured (30s), with retry on failure
