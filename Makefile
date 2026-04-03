.PHONY: check fe-dev be-dev fe-test be-test fe-lint be-lint fe-build be-build migrate-up migrate-down

check: fe-lint be-lint fe-test be-test fe-build be-build

fe-dev:
	cd frontend && pnpm dev

be-dev:
	cd backend && go run ./cmd/api

fe-test:
	cd frontend && pnpm test

be-test:
	cd backend && go test -race ./...

fe-lint:
	cd frontend && pnpm lint && pnpm typecheck

be-lint:
	cd backend && golangci-lint run ./...

fe-build:
	cd frontend && pnpm build

be-build:
	cd backend && go build -o bin/api ./cmd/api

migrate-up:
	cd backend && go run ./cmd/migrate up

migrate-down:
	cd backend && go run ./cmd/migrate down
