.PHONY: build 
build: 
	docker compose build
up:
	docker compose up -d
down:
	docker compose down
logs:
	docker compose logs -f
migrate:
	docker compose exec -it rockbot-db psql -U rockbot -d rockbotdb -f /etc/migration/20260301.00.sql

.PHONY: test
test:
	go test -v -coverprofile=c.out ./...

generate-coverage:
	go tool cover -html=c.out -o coverage.html