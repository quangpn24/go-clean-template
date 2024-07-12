.PHONY: run local-db db/migrate mock lint test testsum

run:
	air -c .air.toml

local-db:
	docker-compose --env-file ./.env -f ./tools/compose/docker-compose.yml -p "go-clean-compose" down
	docker-compose --env-file ./.env -f ./tools/compose/docker-compose.yml -p "go-clean-compose" up -d

db/migrate:
	go run ./cmd/migrate

mock:
	@mockery --name ITransactionUseCase --with-expecter --filename mock_transaction_use_case.go --dir usecase/interfaces --output usecase/mocks
	@mockery --name IPaymentServiceProvider --with-expecter --filename mock_payment_service.go --dir usecase/interfaces --output usecase/mocks
	@mockery --name ITransactionRepository --with-expecter --filename mock_transaction_repo.go --dir usecase/interfaces --output usecase/mocks
	@mockery --name INotifier --with-expecter --filename mock_notifier.go --dir usecase/interfaces --output usecase/mocks
	@mockery --name IDBTransaction --with-expecter --filename mock_db_transaction.go --dir usecase/interfaces --output usecase/mocks

lint:
	@(hash golangci-lint 2>/dev/null || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(go env GOPATH)/bin v1.54.2)
	@golangci-lint run

test:
	go clean -testcache
	go test -cover ./... -gcflags=all=-l -coverprofile  cover.out
	go tool cover -html=cover.out

testsum:
	# go install gotest.tools/gotestsum@latest
	go clean -testcache
	gotestsum --format testname
