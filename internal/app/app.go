package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"booking-service/internal/config"
	httpx "booking-service/internal/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg config.Config
	db  *pgxpool.Pool
	srv *http.Server
}

func New() (*App, error) {
	cfg := config.Load()

	db, err := newDB(cfg.DB.DSN, cfg.Timeouts.DB)
	if err != nil {
		return nil, fmt.Errorf("db init: %w", err)
	}

	router := httpx.NewRouter(db, cfg)

	srv := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &App{cfg: cfg, db: db, srv: srv}, nil
}

func (a *App) Run() error {
	return a.srv.ListenAndServe()
}

func (a *App) Close(ctx context.Context) error {
	if a.db != nil {
		a.db.Close()
	}
	if a.srv != nil {
		return a.srv.Shutdown(ctx)
	}
	return nil
}

func newDB(dsn string, timeout time.Duration) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pcfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
