build:
	@go build -o bin/simpleHolidayServer cmd/main.go

run: build
	@./bin/simpleHolidayServer

test:
	@set -euo pipefail
	@go test -json -v ./... 2>&1 | tee ./logs/gotest.log | gotestfmt

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations/ $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

seeding:
	@go run cmd/seeder/main.go
