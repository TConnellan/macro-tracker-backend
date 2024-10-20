package data

import (
	"context"
	"embed"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed goose_migrations/*.sql
var embedMigrations embed.FS

func newTestDB(t *testing.T) *pgxpool.Pool {

	connConfig, err := pgxpool.ParseConfig("postgresql://test_web:pass@localhost:5432/test_macrotracker")
	if err != nil {
		t.Fatal(err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), connConfig)

	if err != nil {
		t.Fatal(err)
	}

	// run migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	sqldb := stdlib.OpenDBFromPool(db)
	if err := goose.Up(sqldb, "goose_migrations"); err != nil {
		panic(err)
	}
	if err := sqldb.Close(); err != nil {
		panic(err)
	}

	// // run testdata setup
	// script, err := os.ReadFile("./testdata/setup.sql")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ctx, cancel := GetDefaultTimeoutContext()
	// defer cancel()
	// _, err = db.Exec(ctx, string(script))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	t.Cleanup(func() {

		// script, err := os.ReadFile("./testdata/teardown.sql")
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// _, err = db.Exec(ctx, string(script))
		// if err != nil {
		// 	t.Fatal(err)
		// }

		sqldb := stdlib.OpenDBFromPool(db)
		if err := goose.Down(sqldb, "goose_migrations"); err != nil {
			panic(err)
		}
		if err := sqldb.Close(); err != nil {
			panic(err)
		}

		db.Close()
	})

	return db
}
