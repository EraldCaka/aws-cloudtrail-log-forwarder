up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

start:
	docker-compose start

stop:
	docker-compose stop

run:
	@go run cmd/server/main.go

webhook:
	@go run cmd/webhook/main.go
