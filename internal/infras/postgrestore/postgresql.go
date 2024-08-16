package postgrestore

import (
	"fmt"
	"go-clean-template/pkg/config"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Options struct {
	DBName   string
	DBUser   string
	Password string
	Host     string
	Port     string
	SSLMode  bool
}

func ParseFromConfig(c *config.Config) Options {
	return Options{
		DBName:   c.DB.Name,
		DBUser:   c.DB.User,
		Password: c.DB.Pass,
		Host:     c.DB.Host,
		Port:     strconv.Itoa(c.DB.Port),
		SSLMode:  c.DB.EnableSSL,
	}
}

func NewDB(opt Options) (*gorm.DB, error) {
	dsnStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		opt.Host,
		opt.DBUser,
		opt.Password,
		opt.DBName,
		opt.Port,
	)

	db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	db = db.Debug()

	return db, nil
}
