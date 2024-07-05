# Go clean architecture template
This is a template for Go project with clean architecture.

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
