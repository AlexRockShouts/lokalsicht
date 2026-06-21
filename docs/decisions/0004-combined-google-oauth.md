---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Combine Google login and GBP authorization into a single OAuth flow

## Context and Problem Statement

The app requires two OAuth scopes from Google: (1) `openid profile email` for user login/authentication, and (2) `business.manage` for Google Business Profile API access. Having two separate OAuth screens in the onboarding flow would cause significant user drop-off. Each additional click/friction loses 20–30% of users.

## Decision Drivers

- Minimize onboarding friction
- Single consent screen for both scopes
- Handle users who don't have a GBP gracefully

## Considered Options

- Combined OAuth flow (chosen)
- Two separate OAuth flows (rejected)

## Decision Outcome

**Chosen: Single Google OAuth consent screen** requesting `openid profile email business.manage` in one scope parameter, via NextAuth.js.

### Consequences

- Good, because one-click onboarding: login + GBP access in one flow
- Good, because fewer drop-off points in the funnel
- Bad, because users may be hesitant to grant `business.manage` scope when all they want is to log in (mitigated: clear explanation text on the consent screen)
- Bad, because NextAuth.js must handle the expanded scope and pass the resulting access token to Go for GBP API calls (adds complexity to the auth handler)

## Implementation Plan

- **Affected paths**: `frontend/src/app/api/auth/[...nextauth]/route.ts`, `backend/internal/interfaces/http/auth_handler.go`
- **Dependencies**: `next-auth@5`, `golang.org/x/oauth2`
- **Patterns to follow**: NextAuth GoogleProvider with expanded `authorization.params.scope`
- **Patterns to avoid**: Separate login/GBP OAuth flows, storing GBP tokens in the NextAuth session cookie
- **Configuration**: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` in both services

### Verification

- [ ] Single Google consent screen lists both `profile` and `business.manage` scopes
- [ ] GBP access token is obtained during login and stored (encrypted) in DB
- [ ] User without GBP sees a helpful onboarding page ("Create your Google Business Profile in 5 minutes") instead of an error
- [ ] Access token refresh works transparently (user does not see re-auth prompts)
