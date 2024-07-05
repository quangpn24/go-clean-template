.PHONY: run local-db db/migrate mock lint test

run:
	air -c .air.toml

local-db:
	docker-compose --env-file ./.env -f ./tools/compose/docker-compose.yml -p "go-clean-compose" down
	docker-compose --env-file ./.env -f ./tools/compose/docker-compose.yml -p "go-clean-compose" up -d

db/migrate:
	go run ./cmd/migrate

mock:
	@mockery --name ITransactionUseCase --with-expecter --filename mock_transaction_use_case.go --dir usecase --output mocks
	@mockery --name IBankService --with-expecter --filename mock_bank_service.go --dir usecase --output mocks
	@mockery --name ITransactionRepository --with-expecter --filename mock_transaction_repo.go --dir usecase --output mocks
	@mockery --name INotifier --with-expecter --filename mock_notifier.go --dir usecase --output mocks
	@mockery --name IDBTransaction --with-expecter --filename mock_db_transaction.go --dir usecase --output mocks

lint:
	@(hash golangci-lint 2>/dev/null || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(go env GOPATH)/bin v1.54.2)
	@golangci-lint run

test:
	go clean -testcache
	go test -cover ./... -gcflags=all=-l -coverprofile  cover.out
	go tool cover -html=cover.out
