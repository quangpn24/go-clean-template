# Go clean architecture template
This is a template for Go project with clean architecture.
## Folder structure
```
go-clean-template
├── cmd
│   ├── httpserver
│   └── migrate
├── entity/domain/model
├── handler // as controller
│   └── httpserver
│       ├── model
│       ├── options.go
│       ├── server.go
│       └── *_handler.go // handle request to client
├── infras
│   ├── paymentsvc // call API to payment service provider
│   ├── notification // push noti
│   └── postgrestore
│       ├── schema
│       ├── postgresql.go // contains actions to connect DB
│       └── *_repo.go // implement repository interfaces
├── migrations // contains migration files
├── mocks // contains mock files used for unit testing
├── pkg // contains common packages
│   ├── apperror
│   ├── config
│   ├── logger
│   ├── sentry
│   └── validation
├── tools
│   ├── compose
│   └── pre-commit
└── usecase
    ├── interface.go // contains all interfaces
    └── ... // specific usecases
```
## Development

### Init local environment
1. Copy file `.env.example` and rename to `.env`

2. Update env vars to fit your local

3. Start local services
    ```shell
    make local-db
    ```

4. Run the server
    ```shell
    make run
    ```

5. Unit test
    ```shell
    make test
    ```

6. Generate Mock
   ```shell
    make mock
   ```
### Linting

```shell
make lint
```

### Create new migration file

```shell
sql-migrate new -env="development" create-users-table
```

- Result: `Created migration migrations/20230908204301-create-user-table.sql`

Then run migration:
```shell
make db/migrate
```
