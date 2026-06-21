---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Support 5 languages with single primary language per user

## Context and Problem Statement

Switzerland is a four-language country (DE, FR, IT, RM). Additionally, English attracts international users. Each user operates in one primary language. The AI response language is determined by the review text, not the UI language.

## Decision Drivers

- Swiss market: German, French, Italian, Romansh
- International: English
- User sets one primary language; app respects it globally
- AI review response language = language of the original review (not user's UI language)

## Considered Options

- 5 languages, primary per user (chosen)
- German only (rejected)
- All 5 equal priority (rejected — RM is too niche for full translation in MVP)

## Decision Outcome

**Chosen: `next-intl`** with DE/FR/IT as full translations, RM/EN as minimal (navigation, errors, key UI). Each user sets their primary language. The AI response language is derived from the review's text, independently of the UI language.

### Consequences

- Good, because covers the entire Swiss market
- Good, because RM/EN minimal investment protects future expansion
- Good, because `next-intl` provides locale routing (`/de/dashboard`, `/fr/dashboard`, etc.)
- Bad, because 3 full language files means 3x the UI development time for string changes
- Bad, because maintaining 5 language files (even if 2 are minimal) adds overhead

## Implementation Plan

- **Affected paths**: `frontend/src/i18n/{de,fr,it,rm,en}.json`, `frontend/src/app/[locale]/`
- **Dependencies**: `next-intl@3`
- **Patterns to follow**: `next-intl` with `[locale]` segment in App Router, `Accept-Language` header for auto-detect
- **Patterns to avoid**: Client-side language switching without URL change, mixing language detection logic in components
- **Configuration**: `next.config.ts` with locale prefix routing

### Verification

- [ ] UI switches correctly via `/de/dashboard`, `/fr/dashboard`, `/it/dashboard`
- [ ] RM and EN JSON files contain all navigation, error, and key UI strings
- [ ] `Accept-Language` header determines default language on first visit
- [ ] AI response language ≠ UI language (French review → French response, even if UI is German)
