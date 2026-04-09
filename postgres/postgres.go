package postgres

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

// Option — functional option for pool configuration.
type Option func(*pgxpool.Config)

// WithAfterConnect sets a callback that runs once per new connection.
func WithAfterConnect(fn func(ctx context.Context, conn *pgx.Conn) error) Option {
	return func(cfg *pgxpool.Config) {
		cfg.AfterConnect = fn
	}
}

type Config struct {
	User     string `envconfig:"POSTGRES_USER"     required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Port     string `envconfig:"POSTGRES_PORT"     required:"true"`
	Host     string `envconfig:"POSTGRES_HOST"     required:"true"`
	DBName   string `envconfig:"POSTGRES_DB_NAME"  required:"true"`
}

type Pool struct {
	*pgxpool.Pool
}

func New(ctx context.Context, c Config, opts ...Option) (*Pool, error) {
	dsn := fmt.Sprintf("user=%s password=%s port=%s host=%s dbname=%s",
		c.User,
		c.Password,
		c.Port,
		c.Host,
		c.DBName,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	for _, opt := range opts {
		opt(cfg)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	return &Pool{Pool: pool}, nil
}

// Migrate запускает миграции из embedded файлов.
func (p *Pool) Migrate(migrationsFS embed.FS, dir string) error {
	db := stdlib.OpenDBFromPool(p.Pool)
	defer db.Close()

	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	log.Info().Msg("postgres: migrations applied")
	return nil
}

func (p *Pool) Close() {
	p.Pool.Close()
	log.Info().Msg("postgres: closed")
}
