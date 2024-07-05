package main

import (
	"go-clean-template/infras/postgrestore"
	"go-clean-template/pkg/config"
	"go-clean-template/pkg/logger"
	"log"
	"strconv"

	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	applogger, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		applogger.Fatalf("cannot load config: %v\n", err)
	}

	db, err := postgrestore.NewDB(postgrestore.Options{
		DBName:   cfg.DB.Name,
		DBUser:   cfg.DB.User,
		Password: cfg.DB.Pass,
		Host:     cfg.DB.Host,
		Port:     strconv.Itoa(cfg.DB.Port),
		SSLMode:  false,
	})
	if err != nil {
		applogger.Fatalf("cannot connecting to db: %v\n", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	sqlDB, err := db.DB()
	if err != nil {
		applogger.Fatalf("cannot get sql db: %v\n", err)
	}

	total, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		applogger.Fatalf("cannot execute migration: %v\n", err)
	}

	applogger.Infof("applied %d migrations\n", total)
}
