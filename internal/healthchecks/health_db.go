//go:build !no_db

package healthchecks

import (
	"context"

	"microservice/internal/db"
)

// TODO: Check if the health check needs to be expanded

// Base is a very basic healthcheck that pings the databse server and returns an
// error if the connection could not be established.
func Base() error {
	return db.Pool().Ping(context.Background())
}
