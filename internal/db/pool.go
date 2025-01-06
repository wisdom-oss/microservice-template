package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"microservice/internal"
)

// This file contains the connection to the database which is automatically
// initialized on import/app startup

// Pool is not initialized at the app startup and needs to be initiatized by
// calling [Connect].
var Pool *pgxpool.Pool

// Errors which are returned if the database configuration is not in order.
var (
	ErrNoDatabaseUser          = errors.New("no database user configured")
	ErrNoDatabasePassword      = errors.New("no database password configured")
	ErrNoDatabaseHost          = errors.New("no database host configured")
	ErrPoolConfigurationFailed = errors.New("unable to initialize database pool")
	ErrPoolPingFailed          = errors.New("unable to ping database via pool")
)

const pgSqlConnString = `user=%s password=%s host=%s port=%d sslmode=%s database=%s`

func Connect() (err error) {
	slog.Info("initializing database connection")

	config := internal.Configuration

	if !config.IsSet("postgres.host") {
		return ErrNoDatabaseHost
	}
	if !config.IsSet("postgres.user") {
		return ErrNoDatabaseUser
	}
	if !config.IsSet("postgres.password") {
		return ErrNoDatabasePassword
	}

	connectionString := fmt.Sprintf(pgSqlConnString,
		config.GetString("postgres.user"), config.GetString("postgres.password"),
		config.GetString("postgres.host"), config.GetInt("postgres.port"),
		config.GetString("postgres.sslmode"), config.GetString("postgres.database"),
	)
	slog.Debug("generated connection string", "connString", connectionString)

	slog.Debug("initializing database pool with connection string", "connString", connectionString)
	Pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrPoolConfigurationFailed.Error(), err)
	}

	slog.Info("validating database connection")
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("%s: %w", ErrPoolPingFailed.Error(), err)
	}
	return nil
}
