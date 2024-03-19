package main

import (
	"context"
	"expvar"
	"log/slog"
	"os"
	"runtime"
	db "simpleblog/db/sqlc"
	"simpleblog/util"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	version = util.Version()
)

type application struct {
	config util.Config
	logger *slog.Logger
	store  db.Store
	wg     sync.WaitGroup
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	//Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//Configuration
	cfg, err := util.LoadConfig(".")
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("configuration loaded")

	//Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//DB connection
	dbpool, err := pgxpool.New(ctx, cfg.DBSource)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer dbpool.Close()

	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() any {
		return dbpool.Stat()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	logger.Info("pgx connection pool established")

	// App
	app := &application{
		config: cfg,
		logger: logger,
		store:  db.NewStore(dbpool),
	}

	if err = app.serve(); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
