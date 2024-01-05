include .env

sqlc-generate:
	@sqlc generate

createmigration:
	@migrate create -ext sql -dir  database/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

pushmigrations:
	@migrate -source file://database/migrations -database $(DATABASE_URL) up

dropmigrations:
	@migrate -source file://database/migrations -database $(DATABASE_URL) down

run-dev:
	@air -c .air.toml

format:
	@golangci-lint run --fix

%:
	@:
