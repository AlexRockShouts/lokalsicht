.PHONY: dev setup test lint build clean

dev: dev-db
	@echo "==> Backend :5174 | Frontend :3000"
	@make -j2 dev-backend dev-frontend

dev-db:
	docker compose up -d postgres
	@sleep 2

dev-backend:
	cd backend && air

dev-frontend:
	cd frontend && npm run dev

setup:
	cp -n .env.example .env || true
	docker compose up -d postgres
	cd backend && go mod download && go run ./cmd/server
	cd frontend && npm install

test:
	cd backend && go test ./... -v
	cd frontend && npm run test -- --run

lint:
	cd backend && golangci-lint run ./...
	cd frontend && npm run lint

typecheck:
	cd frontend && npx tsc --noEmit

swagger:
	cd backend && swag init -g cmd/server/main.go -o docs/

typesync: swagger
	cd frontend && npx openapi-typescript ../backend/docs/swagger.json -o src/api/types.ts

build:
	cd backend && CGO_ENABLED=0 go build -o server ./cmd/server
	cd frontend && npm run build

docker-down:
	docker compose down

clean:
	docker compose down -v
	rm -f backend/server
	rm -rf frontend/.next
