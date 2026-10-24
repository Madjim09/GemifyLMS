.PHONY: up down stop logs build ps reset db test

up:
	docker compose up --build

down:
	docker compose down

stop:
	docker compose stop

logs:
	docker compose logs -f

build:
	docker compose build

ps:
	docker compose ps

reset:
	docker compose down -v

db:
	docker compose up -d postgres redis

test:
	cd backend && go test ./...