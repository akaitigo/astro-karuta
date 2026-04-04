# ADR-003: Deferred Authentication

## Status
Accepted

## Context
Astro-Karuta is currently in MVP phase. Full user authentication (login, session
management, token validation) adds significant complexity that is not required
for the core educational game loop.

However, several security-sensitive paths exist:

1. **REST endpoints** (`/api/v1/collections`, `/api/v1/missions`) accept a
   `user_id` query parameter without verifying ownership.
2. **WebSocket reconnect** (`HandleReconnect`) allows a new client to resume
   any player session by providing the correct `game_id` and `player_id`.
3. **Mission completion** credits bonus cards to any supplied `user_id`.

## Decision
- **MVP (current):** Validate `user_id` format (UUID) at the handler level.
  Reject empty or malformed values with 400 Bad Request. Log warnings when
  WebSocket reconnect occurs from a different `clientID` than the original.
- **Post-MVP:** Implement session-based or JWT authentication.
  - REST endpoints: require `Authorization` header; derive `user_id` from
    the verified token.
  - WebSocket reconnect: require a reconnect token issued at initial
    connection. Reject reconnect attempts without a valid token.

## Consequences
- MVP users can impersonate other users by guessing UUIDs. This is acceptable
  because the MVP is not publicly deployed and UUIDs are random.
- Post-MVP must introduce authentication before any public deployment.
- The reconnect warning log enables monitoring for suspicious activity
  during the MVP period.
