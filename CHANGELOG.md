# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2026-04-04

### Features

- **Card/Deck CRUD API with seed data** -- 49 astronomy cards (constellations, planets, phenomena) with category filtering, deck management, and seasonal deck endpoint (#6)
- **WebSocket battle engine with matchmaking** -- Real-time 2-player karuta matches via WebSocket; room-based and random matchmaking; server-side grab judgment with timestamp comparison; reconnection support within 30 seconds (#7)
- **Collection feature** -- Backend API and frontend UI for tracking collected cards; collection progress stats with percentage; category-based filtering (#8)
- **Seasonal deck and observation missions** -- Season-aware deck generation based on current date; observation missions with GPS-based completion validation; bonus card rewards (#9)
- **Game frontend UI** -- Lobby page with room creation/join/random match; real-time battle interface with reading cards and picture cards; score board; game result screen with winner announcement (#10)

### Bug Fixes

- **Security and robustness fixes** -- Input validation (room code length, user name sanitization), WebSocket origin check, rate limiting, request timeout middleware (#12, #13)
- **errcheck lint and UUID user ID** -- All error returns properly checked; user IDs changed from sequential int to UUID; Promise leak in useGame hook fixed (#14)
- **Game timeout and reconnect state** -- Game timeout after 10 minutes inactivity; reconnect restores full game state; graceful shutdown with in-flight game drain; matchmaking queue cleanup on disconnect (#15)

### Infrastructure

- **Monorepo scaffold** -- Next.js 15 frontend (PWA) + Go 1.23 backend in a single repository with Makefile commands
- **CI pipeline** -- GitHub Actions workflow for both frontend (lint, typecheck, test, build) and backend (golangci-lint, test with race detector, build)
- **Architecture Decision Records** -- ADR-001 (WebSocket over SSE), ADR-002 (server-side timestamp for grab judgment)
