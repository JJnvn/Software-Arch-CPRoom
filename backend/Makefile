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

.PHONY: restart-% build-%
restart-%:
	docker-compose build $*
	docker-compose up -d $*

.PHONY: restart-approval restart-booking restart-auth restart-room restart-notification
restart-approval: restart-approval-service
restart-booking: restart-booking-service
restart-auth: restart-auth-service
restart-room: restart-room-service
restart-notification: restart-notification-service
