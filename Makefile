SHELL=cmd

# run up migrations, user details based on docker-compose.yml
migrate_up:
	migrate -path db/migrations -database "postgresql://devuser:admin@localhost:5432/devdb?sslmode=disable" -verbose up

# run down migrations, user details based on docker-compose.yml
migrate_down:
	migrate -path db/migrations -database "postgresql://devuser:admin@localhost:5432/devdb?sslmode=disable" -verbose down

# generate db related go code with sqlc
sqlc:
	docker run --rm -v "${PWD}:/src" -w /src kjconroy/sqlc generate

# run all tests
test:
	go test -v -cover ./...

.PHONY: migrate_up, migrate_down, sqlc, test