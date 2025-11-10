.PHONY: reset up down
reset:
	docker compose down -v
	docker compose up -d --build

up:
	docker compose up -d --build

down:
	docker compose down -v

.PHONY: loadtest
loadtest:
	go run loadtest.go