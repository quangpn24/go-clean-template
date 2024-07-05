package main

import (
	"fmt"
	"go-clean-template/handler/httpserver"
	"go-clean-template/infras/banksv"
	"go-clean-template/infras/notification"
	"go-clean-template/infras/postgrestore"
	"go-clean-template/pkg/config"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/sentry"
	"go-clean-template/usecase"
	"log"

	sentrygo "github.com/getsentry/sentry-go"
)

func main() {
	applog, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}
	defer logger.Sync(applog)

	cfg, err := config.LoadConfig()
	if err != nil {
		applog.Fatal(err)
	}

	err = sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.AppEnv,
		AttachStacktrace: true,
	})
	if err != nil {
		applog.Fatalf("cannot init sentry: %v", err)
	}
	defer sentrygo.Flush(sentry.FlushTime)

	db, err := postgrestore.NewDB(postgrestore.ParseFromConfig(cfg))
	if err != nil {
		applog.Fatal(err)
	}

	server, err := httpserver.New(httpserver.WithConfig(cfg), httpserver.WithLogger(applog))
	if err != nil {
		applog.Fatal(err)
	}

	//Setup Dependencies
	transRepo := postgrestore.NewTransactionRepo(db)
	bankSvc := banksv.NewBankService()
	dbTransaction := postgrestore.NewDBTransaction(db)
	transUseCase := usecase.NewTransactionUseCase(transRepo, bankSvc, dbTransaction)
	transUseCase.SetNotifiers(notification.NewEmailNotifier(), notification.NewAppNotifier())

	server.TransactionUseCase = transUseCase

	addr := fmt.Sprintf(":%d", cfg.Port)
	applog.Fatal(server.Start(addr))
}
