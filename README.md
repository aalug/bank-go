# Go Bank

### App built in Go 1.20

## The app uses:
- Postgres
- Docker
- [Gin](https://github.com/gin-gonic/gin)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [testify](https://github.com/stretchr/testify)
- [PASETO Security Tokens](github.com/o1egl/paseto)
- [jwt-go](https://github.com/dgrijalva/jwt-go)

## Getting started
1. Clone the repository
2. Go to the project's root directory
3. Rename `app.env.sample` to `app.env` and replace the values
4. Run in your terminal:
    - `docker-compose up --build` to run the containers
5. Now everything should be ready and the server running on `SERVER_ADDRESS` specified in app.env

## Testing
1. Run the containers (database container is sufficient)
2. Run in your terminal:
    - `make test`
   or
    - `make test_coverage p={PATH}` - to get the coverage in the HTML format - where `{PATH}` is the path to the target directory for which you want to generate test coverage. The `{PATH}` should be replaced with the actual path you want to use. For example `./api`
   or
    - simply use `go test` commands

## API Endpoints
#### Users (only endpoints that do not require authentication)
 - `/users` - handles POST requests to create users
 - `/users/login` - handles POST requests to log in users
 - `/tokens/renew` - handles  POST requests to renew the access tokens

### Accounts
- `/accounts` - handles POST requests to create accounts
- `/accounts` - handles GET requests to get all accounts
- `/accounts/{id}` - handles GET requests to get account details
- `/accounts/{id}` - handles DELETE requests to delete an account

### Transfers
- `/transfers` - handles POST requests to transfer money from one account to another


## Database
The database's schema and intricate details can be found on 
dedicated webpage, which provides a comprehensive overview 
of the data structure, tables, relationships, and other essential 
information. To explore the database further, please visit
this [dbdocs.io webpage](https://dbdocs.io/aalug/bank_go).
Password: `bankgopassword`
