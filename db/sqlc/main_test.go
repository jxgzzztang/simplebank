package db

import (
	"context"
	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	dbSource = "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

var testQuery *Queries

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error

	testDB, err = pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("error opening db:", err)
	}

	testQuery = New(testDB)

	defer testDB.Close()

	os.Exit(m.Run())

}
