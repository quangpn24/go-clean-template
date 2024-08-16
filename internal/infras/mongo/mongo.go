package mongo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-clean-template/pkg/config"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Options struct {
	DBName string
	Host   string
	Port   string
	User   string
	Pass   string
}

func ParseFromConfig(c *config.Config) Options {
	return Options{
		DBName: c.MongoDB.DBName,
		Host:   c.MongoDB.Host,
		Port:   strconv.Itoa(c.MongoDB.Port),
		User:   c.MongoDB.User,
		Pass:   c.MongoDB.Pass,
	}
}

func NewDB(opt Options) (*mongo.Database, error) {
	uri := fmt.Sprintf("mongodb://%v:%v", opt.Host, opt.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			fmt.Println(e.Command)
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			fmt.Println(e.Reply)
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			fmt.Println(e.Failure)
		},
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: opt.User,
		Password: opt.Pass,
	}).SetMonitor(monitor))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(opt.DBName), nil
}
