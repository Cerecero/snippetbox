package models 

import (
	"database/sql"
	"os"
	"testing"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func newTestDB(t *testing.T) *sql.DB {
	// testDB := os.Getenv("TESTDB_URL")

	// db, err := sql.Open("pgx", testDB)
	db, err := sql.Open("pgx", "postgres://test_web:pass@/test_snippetbox")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}
