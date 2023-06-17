# Go Bank

### App built in Go 1.20

## The app uses:
- Postgres
- Docker
- [Gin](https://github.com/gin-gonic/gin)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [testify](https://github.com/stretchr/testify)
- [jwt-go](https://github.com/dgrijalva/jwt-go)
- [PASETO Security Tokens](github.com/o1egl/paseto)

## Getting started
1. Clone the repository
2. Go to the project's root directory
3. Run in your terminal:
    - `docker-compose up` to run the database container
    - `make migrate_up` to run migrations

## Testing
1. Run the postgres container (`docker-compose up`)
2. Run in your terminal:
    - `make test`
   or
    - `make test_coverage`