# Harvest Report: Astro-Karuta v1.0.0

**Date**: 2026-04-04
**Tag**: v1.0.0
**Repository**: akaitigo/astro-karuta

---

## Project Summary

Astro-Karuta is a real-time multiplayer educational game that teaches astronomy through the traditional Japanese card game "karuta". Players compete to grab the correct picture card when a reading clue is announced via WebSocket, collecting celestial objects (constellations, planets, phenomena) along the way.

### Tech Stack

| Layer     | Technology                             |
|-----------|----------------------------------------|
| Frontend  | TypeScript, Next.js 15, React 19, PWA |
| Real-time | WebSocket (gorilla/websocket)          |
| Backend   | Go 1.23, net/http stdlib router        |
| Testing   | Vitest + Testing Library, go test      |
| CI        | GitHub Actions                         |

---

## Statistics

### Issues

| Metric     | Count |
|------------|-------|
| Total      | 5     |
| Closed     | 5     |
| Open       | 0     |

All 5 MVP issues were completed and closed:
1. F1: Card/Deck CRUD API (#1)
2. F1: WebSocket battle engine (#2)
3. F2: Collection feature (#3)
4. F1: Game frontend UI (#4)
5. F3: Seasonal deck and observation missions (#5)

### Pull Requests

| Metric | Count |
|--------|-------|
| Total  | 7     |
| Merged | 7     |
| Open   | 0     |

| PR # | Title | Type |
|------|-------|------|
| #6   | Card/Deck CRUD API with seed data | feat |
| #7   | WebSocket battle engine with matchmaking | feat |
| #8   | Collection feature (backend + frontend) | feat |
| #9   | Seasonal deck and observation missions | feat |
| #10  | Game frontend UI | feat |
| #11  | README and CHANGELOG for v1.0.0 release | docs |
| #12  | Review round 1 - security and robustness | fix |

**Merge rate**: 7/7 (100%)

### Tests

| Scope    | Files | Tests |
|----------|-------|-------|
| Backend  | 6     | 76    |
| Frontend | 10    | 65    |
| **Total**| **16**| **141**|

All tests pass with race detector enabled (backend) and no lint/typecheck errors.

### Codebase Size

| Directory | Go files | TS/TSX files |
|-----------|----------|--------------|
| backend/  | 31       | -            |
| frontend/ | -        | 34           |

---

## Review Loop History

### Round 1

**Findings**: 13 total (2 CRITICAL, 5 HIGH, 4 MEDIUM, 2 LOW)

| # | Severity | Category | File | Issue | Resolution |
|---|----------|----------|------|-------|------------|
| 1 | CRITICAL | Security | ws_handler.go | WebSocket CheckOrigin accepted all origins | Fixed: validate against CORS_ORIGIN env var |
| 2 | CRITICAL | Security | config.go | Hardcoded DB credentials in default config | Fixed: empty defaults, require env var |
| 3 | HIGH | Robustness | game_manager.go | createGameState nil return not checked | Fixed: nil checks + error messages |
| 4 | HIGH | Robustness | game_manager.go | MarshalMessage errors ignored (10 sites) | Fixed: all error returns handled with logging |
| 5 | HIGH | Correctness | game_manager.go | Incorrect grab advances card index | Fixed: only correct grabs advance |
| 6 | HIGH | Security | ws_handler.go | No connection limit on WebSocket | Fixed: 200 max concurrent connections |
| 7 | HIGH | Dead Code | config.go | Config unused in main.go | Noted: config ready for DB integration phase |
| 8 | MEDIUM | Consistency | lobby/game pages | WS URL default path mismatch | Fixed: /ws -> /api/v1/ws |
| 9 | MEDIUM | Consistency | missions/page.tsx | User ID constant differs from collection page | Fixed: unified to DEFAULT_USER_ID |
| 10 | MEDIUM | Robustness | game_manager.go | Fewer than 4 candidates possible | Deferred: edge case with <4 cards in deck |
| 11 | LOW | Config | config.go | Unused Config struct | Not fixed: ready for DB phase |

**Result**: CRITICAL=0, HIGH=0 after Round 1. Loop terminated.

---

## Lessons Learned

1. **WebSocket security is easy to overlook** -- `CheckOrigin: return true` was the first thing written for development convenience and never updated. Origin validation should be the default from the start, with an explicit opt-out for local dev only.

2. **Error handling on marshaling** -- JSON marshaling errors for well-known structs feel "impossible" but ignoring them with `_` creates silent failures. A helper function that logs-and-returns-nil is better than scattered `_` assignments.

3. **Game logic correctness matters early** -- The incorrect-grab-advances-card bug would have been caught by a game-play integration test but was missed because tests only verified the message type, not the game state progression. Test assertions should verify state changes, not just response shapes.

4. **Hardcoded credentials in defaults** -- Even "development" defaults like `postgres://user:password@...` get committed and stay forever. Use empty defaults and fail fast.

5. **Frontend-backend URL consistency** -- WebSocket paths drifted between frontend defaults and backend routes. A shared constants file or `.env.example` that documents all endpoints would prevent this.

---

## Template Improvement Suggestions

1. **Add security checklist to PR template** -- Include items for origin validation, input size limits, and credential scanning before merge.

2. **Require error-handling linter** -- Enable `errcheck` in golangci-lint from day one. The 10 ignored MarshalMessage errors would have been caught automatically.

3. **Add game-state integration tests** -- The test template should include a multi-turn game simulation that verifies state transitions, not just message types.

4. **Standardize frontend config** -- Create a shared `config.ts` that centralizes API_BASE and WS_URL with clear documentation, rather than scattering `process.env` lookups across pages.

5. **Security review as Stage 4 default** -- Move security review earlier in the pipeline (e.g., pre-merge hook or CI check for `CheckOrigin.*return true` patterns).
