.PHONY: reset
reset:
	docker compose down -v
	docker compose up -d --build


.PHONY: loadtest
loadtest:
	go run loadtest.go