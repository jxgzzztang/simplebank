package db

import (
	"context"
	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jxgzzztang/simplebank/util"
	"log"
	"os"
	"testing"
)


var testQuery *Queries

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error

	err = util.LoadConfig("../..")

	testDB, err = pgxpool.New(ctx, util.Config.DBSource)
	if err != nil {
		log.Fatal("error opening db:", err)
	}

	testQuery = New(testDB)

	defer testDB.Close()

	os.Exit(m.Run())

}
