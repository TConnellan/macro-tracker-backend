package data

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed goose_migrations/*.sql
var embedMigrations embed.FS

//go:embed testdata/setup.sql
var setupSQL string

// go: embed testdata/teardown.sql
var teardownSQL string

func newTestDB(t *testing.T, testName string) (*pgxpool.Pool, error) {

	defaultDBDSN := os.Getenv("DEFAULT_DB_DSN")
	testDBDSN := os.Getenv("TEST_DB_DSN")
	testUser := os.Getenv("TEST_USER")
	testPass := os.Getenv("TEST_PASS")

	defaultConnConfig, err := pgxpool.ParseConfig(defaultDBDSN)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	defaultDB, err := pgxpool.NewWithConfig(ctx, defaultConnConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer defaultDB.Close()

	// Construct the unique database name
	dbName := fmt.Sprintf("test_macrotracker_%s_%s", testName, uuid.NewString())

	// Step 1: Create the database
	_, err = defaultDB.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s WITH OWNER 'postgres' ENCODING 'UTF8' TABLESPACE 'pg_default';`, pq.QuoteIdentifier(dbName)))
	if err != nil {
		t.Fatal(fmt.Errorf("failed to create database: %w", err))
	}

	// create the user
	_, err = defaultDB.Exec(ctx, fmt.Sprintf(`CREATE USER %s WITH PASSWORD '%s';`, pq.QuoteIdentifier(testUser), pq.QuoteIdentifier(testPass)))
	if err != nil {
		if !strings.Contains(err.Error(), `role "`+testUser+`" already exists`) {
			t.Fatal(fmt.Errorf("error when creating user: %w", err))
		}
	}

	testConnConfig, err := pgxpool.ParseConfig(fmt.Sprintf(testDBDSN, dbName))
	if err != nil {
		t.Fatal(err)
	}

	testDB, err := pgxpool.NewWithConfig(ctx, testConnConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Step 3: Execute parameterized queries within the new database
	_, err = testDB.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO %s; GRANT CREATE ON SCHEMA public TO %s; CREATE EXTENSION IF NOT EXISTS citext;`, testUser, testUser))
	if err != nil {
		t.Fatal(fmt.Errorf("error when granting user permissions to %s: %w", dbName, err))
	}

	_, err = testDB.Exec(ctx, fmt.Sprintf(`GRANT CREATE, CONNECT ON DATABASE %s TO integration_tester;`, pq.QuoteIdentifier(dbName)))
	if err != nil {
		t.Fatal(fmt.Errorf("failed to grant privileges to %s: %w", dbName, err))
	}

	// run migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatal(err)
	}

	sqldb := stdlib.OpenDBFromPool(testDB)
	if err := goose.Up(sqldb, "goose_migrations"); err != nil {
		t.Fatal(err)
	}
	if err := sqldb.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = testDB.Exec(ctx, setupSQL)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer testDB.Close()

		_, err = testDB.Exec(context.Background(), teardownSQL)
		if err != nil {
			t.Fatal(err)
		}

		sqldb := stdlib.OpenDBFromPool(testDB)
		defer func() {
			if err := sqldb.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		if err := goose.Down(sqldb, "goose_migrations"); err != nil {
			t.Fatal(err)
		}

	})

	return testDB, nil
}
