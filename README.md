# Astro-Karuta

<!-- badges -->
[![CI](https://github.com/akaitigo/astro-karuta/actions/workflows/ci.yml/badge.svg)](https://github.com/akaitigo/astro-karuta/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> Learn astronomy through the traditional Japanese card game "Karuta".

Astro-Karuta is a real-time multiplayer educational game that teaches children (ages 6-15) about constellations, planets, and astronomical phenomena through digital karuta matches. Players compete to grab the correct picture card when a reading clue is announced, collecting celestial objects along the way.

---

## Quick Start

### Prerequisites

| Tool   | Version |
|--------|---------|
| Go     | 1.23+   |
| Node.js| 22+     |
| pnpm   | 10+     |

### Clone & Setup

```bash
git clone git@github.com:akaitigo/astro-karuta.git
cd astro-karuta

# Backend
cd backend
cp .env.example .env
go mod download
cd ..

# Frontend
cd frontend
cp .env.example .env.local
pnpm install
cd ..
```

### Run

```bash
# Start backend (default :8080)
make be-dev

# Start frontend (default :3000)
make fe-dev
```

### Test & Build

```bash
# Run everything: lint + test + build
make check

# Or individually
make be-test      # Go tests with race detector
make fe-test      # Vitest
make fe-lint      # ESLint + TypeScript check
make be-build     # Compile backend binary
make fe-build     # Next.js production build
```

---

## Tech Stack

| Layer     | Technology                                  |
|-----------|---------------------------------------------|
| Frontend  | TypeScript, Next.js 15 (PWA), React 19      |
| Real-time | WebSocket (gorilla/websocket)               |
| Backend   | Go 1.23, net/http (stdlib router)           |
| Testing   | Vitest + Testing Library (FE), go test (BE) |
| CI        | GitHub Actions                              |
| Data      | In-memory (no external DB required)         |

---

## Architecture

```
┌──────────────────┐         WebSocket          ┌──────────────────┐
│                  │◄──────────────────────────► │                  │
│  Next.js PWA     │         /api/v1/ws          │  Go API Server   │
│  (port 3000)     │                             │  (port 8080)     │
│                  │────── REST /api/v1/* ──────► │                  │
└──────────────────┘                             └──────────────────┘
        │                                                │
        │  React Components                              │  Internal Layers
        │  ├── Lobby (matchmaking)                       │  ├── handler/
        │  ├── Game (battle UI)                          │  ├── service/
        │  ├── Collection (card gallery)                 │  ├── repository/
        │  ├── Seasonal (deck browser)                   │  ├── ws/ (hub + game manager)
        │  └── Missions (observation)                    │  ├── model/
        │                                                │  └── pkg/astronomy/
        │                                                │
        └── hooks/                                       └── seed/ (49 cards)
            ├── useWebSocket (reconnect + backoff)
            └── useGame (state machine)
```

### Game Flow

1. **Lobby** -- Player enters name, chooses room code / random match
2. **Matchmaking** -- WebSocket `join` message; server pairs two players
3. **Battle** -- Server sends `card_revealed` with reading text + candidate cards; players race to `grab` the correct one
4. **Judgment** -- Server compares timestamps, broadcasts `grab_result`
5. **Game Over** -- After all cards are played, server sends `game_over` with final scores
6. **Collection** -- Won cards are added to the player's collection

---

## API Endpoints

### REST

| Method | Path                             | Description                  |
|--------|----------------------------------|------------------------------|
| GET    | `/api/v1/health`                 | Health check                 |
| GET    | `/api/v1/cards`                  | List cards (filter: category, season) |
| GET    | `/api/v1/cards/{id}`             | Get card detail              |
| GET    | `/api/v1/decks`                  | List all decks               |
| GET    | `/api/v1/decks/{id}`             | Get deck by ID               |
| GET    | `/api/v1/decks/seasonal`         | Get current seasonal deck    |
| GET    | `/api/v1/collections`            | List user collection (query: user_id, category) |
| GET    | `/api/v1/collections/stats`      | Collection progress stats    |
| GET    | `/api/v1/missions`               | List active missions (query: user_id) |
| POST   | `/api/v1/missions/{id}/complete` | Complete a mission (body: user_id, lat, lng) |

### WebSocket

| Direction | Message Type    | Description                     |
|-----------|-----------------|---------------------------------|
| Client    | `join`          | Join/create room or random match|
| Client    | `grab`          | Attempt to grab a card          |
| Client    | `reconnect`     | Reconnect to existing game      |
| Server    | `player_joined` | Notify new player in room       |
| Server    | `player_left`   | Notify player disconnection     |
| Server    | `card_revealed` | New card reading + candidates   |
| Server    | `grab_result`   | Result of grab attempt          |
| Server    | `game_over`     | Final scores and winner         |
| Server    | `match_found`   | Random match paired             |
| Server    | `waiting`       | Waiting for opponent             |
| Server    | `error`         | Error message                   |

---

## Project Structure

```
astro-karuta/
├── backend/
│   ├── cmd/api/          # Entry point
│   ├── config/           # Configuration
│   ├── internal/
│   │   ├── handler/      # HTTP + WS handlers
│   │   ├── middleware/   # CORS, rate limiting, timeout
│   │   ├── model/        # Domain models
│   │   ├── repository/   # Data access (in-memory)
│   │   ├── seed/         # Card seed data (49 cards)
│   │   ├── service/      # Business logic
│   │   └── ws/           # WebSocket hub + game manager
│   ├── migrations/       # SQL migrations
│   └── pkg/astronomy/    # Seasonal calculation
├── frontend/
│   └── src/
│       ├── app/          # Next.js pages (lobby, game, collection, seasonal, missions)
│       ├── components/   # React components
│       ├── hooks/        # useWebSocket, useGame, useAudio
│       ├── lib/          # API client, test setup
│       └── types/        # TypeScript type definitions
├── docs/
│   ├── PRD.md            # Product Requirements Document
│   └── adr/              # Architecture Decision Records
├── .github/workflows/    # CI pipeline
└── Makefile              # Build commands
```

---

## License

MIT
