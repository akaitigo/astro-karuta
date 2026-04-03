# Astro-Karuta

天文知識をかるた形式で学ぶ子ども向け教育ゲーム。

## Tech Stack
- Frontend: TypeScript / Next.js (PWA), Web Audio API, WebSocket
- Backend: Go, PostgreSQL, Redis
- Infra: GCP Cloud Run, Cloud Storage

## Structure
```
frontend/    — Next.js PWA (pnpm)
backend/     — Go API + WebSocket server
docs/        — PRD, ADR, harvest
migrations/  — SQL migration files (in backend/)
```

## Commands
```bash
make check          # lint + test + build (全体)
make fe-dev         # frontend dev server
make be-dev         # backend dev server
make fe-test        # frontend tests
make be-test        # backend tests
make fe-lint        # frontend lint + typecheck
make be-lint        # backend lint (golangci-lint)
make migrate-up     # DB migration 適用
make migrate-down   # DB migration ロールバック
```

## Rules
- Go: エラーは必ずハンドリング、context.Context第1引数
- TypeScript: `any` 禁止、`as` 最小限
- SQL: パラメータバインド必須
- テスト: 正常系+異常系を必ず書く
- シークレット: 環境変数管理、コードに含めない
