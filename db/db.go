package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB represents a database
type DB struct {
}

// Database interface allows a simple test for a connection to the DB
type Database interface {
	TestConnection(ctx context.Context) error
}

// TestConnection ensures a connection to the DB
func (*DB) TestConnection(ctx context.Context) error {
	pool, err := pgxpool.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	pool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		err := conn.Ping(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
