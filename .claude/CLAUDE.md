# Internal Dev Notes

## Architecture
- Monorepo: frontend/ (Next.js) + backend/ (Go)
- WebSocket: backend handles real-time game state
- Match flow: lobby → matchmaking → game → result → collection update

## Conventions
- Frontend components: PascalCase, co-located tests
- Backend packages: internal/ for private, pkg/ for shared logic
- API routes: /api/v1/ prefix
- DB: snake_case for columns, migrations in backend/migrations/
- Env: .env.example in both frontend/ and backend/

## Testing
- Frontend: vitest + @testing-library/react
- Backend: go test with testcontainers for DB tests
- WebSocket: integration tests with gorilla/websocket client

## PR Workflow
1. Branch: feat/, fix/, docs/ prefix
2. make check must pass before PR
3. Squash merge to main
