
# run up migrations, user details based on docker-compose.yml
migrate_up:
	migrate -path db/migrations -database "postgresql://devuser:admin@localhost:5432/devdb?sslmode=disable" -verbose up

# run down migrations, user details based on docker-compose.yml
migrate_down:
	migrate -path db/migrations -database "postgresql://devuser:admin@localhost:5432/devdb?sslmode=disable" -verbose down

# generate db related go code with sqlc
sqlc:
	cmd.exe /c "docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate"

# generate mock db for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/aalug/go-bank/db/sqlc Store

# run all tests
test:
	go test -v -cover ./...

# run tests in the given path (p) and display results in the html file
test_coverage:
	go test $(p) -coverprofile=coverage.out && go tool cover -html=coverage.out

# run the HTTP server
runserver:
	go run main.go

.PHONY: migrate_up, migrate_down, sqlc, test, test_coverage, runserver, mock