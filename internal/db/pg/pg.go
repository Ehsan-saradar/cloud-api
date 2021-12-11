package pg

import (
	"context"
	"database/sql"
	"fmt"

	"api.cloud.io/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	Query func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	TheDB *sqlx.DB
)

type Client struct {
	db            *sqlx.DB
	logger        zerolog.Logger
	migrationsDir string
}

func NewClient(cfg config.PgConfig) (*Client, error) {
	err := createDB(cfg.Host, cfg.Port, cfg.Sslmode, cfg.UserName, cfg.Password, cfg.Database)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create database %s", cfg.Database)
	}
	logger := log.With().Str("module", "pg").Logger()
	db, err := openDB(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not open database connection")
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	Query = db.QueryContext
	TheDB = db
	cli := &Client{
		db:            db,
		logger:        logger,
		migrationsDir: cfg.MigrationsDir,
	}
	if err := cli.Migrate(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to run migrations")
	}
	er := initGames(cfg.JsonDir)
	return cli, er
}

func createDB(host string, port int, ssl, username, password, name string) error {
	connStr := fmt.Sprintf("user=%s sslmode=%v password=%v host=%v port=%v", username, ssl, password, host, port)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return errors.Wrap(err, "failed to open postgres connection")
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%v');`, name)
	row := db.QueryRow(query)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if !exists {
		query = fmt.Sprintf(`CREATE DATABASE %v`, name)
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func openDB(cfg config.PgConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%v host=%v port=%v", cfg.UserName, cfg.Database, cfg.Sslmode, cfg.Password, cfg.Host, cfg.Port)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return &sqlx.DB{}, err
	}

	return db, nil
}

func (s *Client) Migrate(cfg config.PgConfig) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%v:%v/%s?sslmode=%v", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Sslmode)
	pgMigrate, err := migrate.New(cfg.MigrationsDir, connStr)
	// pgMigrate, err := migrate.New(cfg.MigrationsDir, connStr)
	if err != nil {
		return err
	}

	err = pgMigrate.Migrate(uint(cfg.MigrateVersion))
	if err != nil && err.Error() != "no change" {
		return err
	}
	s.logger.Debug().Int("Applied migrations", cfg.MigrateVersion)
	return nil
}
