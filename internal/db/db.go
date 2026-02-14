package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

var (
	CreateExtensionSQL = "CREATE EXTENSION IF NOT EXISTS vector"

	SetupTableSQL = `CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title TEXT,
		overview TEXT,
		genres TEXT[],
		language TEXT,
		popularity FLOAT8,
		release_date DATE,
		embedding vector(3072),
		embedded_at TIMESTAMPTZ
	);`
	DropTableSQL = "DROP TABLE IF EXISTS movies"
)

// OpenDB opens a connection to the database and returns the connection along with a close function.
func OpenDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, func(), error) {
	pgcfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	pgcfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, CreateExtensionSQL)
		if err != nil {
			return fmt.Errorf("unable to create vector extension: %v", err)
		}
		return pgxvec.RegisterTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgcfg)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("unable to ping database: %v", err)
	}

	slog.InfoContext(ctx, "Setting up movies table...")
	_, err = pool.Exec(ctx, SetupTableSQL)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create table: %v", err)
	}

	slog.Info("Connected to database")
	return pool, pool.Close, nil
}
