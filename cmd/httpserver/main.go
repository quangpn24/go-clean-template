package main

import (
	"fmt"
	"log"

	"go-clean-template/internal/handler/httpserver"
	"go-clean-template/internal/infras/mongo"
	"go-clean-template/internal/infras/notification"
	"go-clean-template/internal/infras/paymentsvc"
	"go-clean-template/internal/usecase"
	"go-clean-template/pkg/config"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/sentry"

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

	//db, err := postgrestore.NewDB(postgrestore.ParseFromConfig(cfg))
	//if err != nil {
	//	applog.Fatal(err)
	//}

	db, err := mongo.NewDB(mongo.ParseFromConfig(cfg))
	if err != nil {
		applog.Fatal(err)
	}

	server, err := httpserver.New(httpserver.WithConfig(cfg), httpserver.WithLogger(applog))
	if err != nil {
		applog.Fatal(err)
	}

	//Setup Dependencies
	//transRepo := postgrestore.NewTransactionRepo(db)
	transRepo := mongo.NewTransactionRepo(db)
	paymentSvc := paymentsvc.NewPaymentServiceProvider()
	transUseCase := usecase.NewTransactionUseCase(transRepo, paymentSvc)
	transUseCase.SetNotifiers(notification.NewEmailNotifier(), notification.NewAppNotifier())

	server.TransactionUseCase = transUseCase

	addr := fmt.Sprintf(":%d", cfg.Port)
	applog.Fatal(server.Start(addr))
}
