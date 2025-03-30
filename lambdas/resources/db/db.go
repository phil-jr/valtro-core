package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func init() {
	ctx := context.Background()

	dsn := "postgresql://postgres.pixydqwfhsdbzbxehilj:txM54gw@yXHixcK@aws-0-us-east-1.pooler.supabase.com:6543/postgres?simple_protocol=true"
	if dsn == "" {
		log.Fatal("Environment variable POSTGRES_DSN is required")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Printf("Failed to parse Postgres DSN: %v", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Printf("Failed to create connection pool: %v", err)
	}

	log.Println("Connection pool initialized")
}
