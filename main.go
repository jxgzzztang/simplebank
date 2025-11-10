package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jxgzzztang/simplebank/api"
	db "github.com/jxgzzztang/simplebank/db/sqlc"
	"github.com/jxgzzztang/simplebank/util"
)

func main() {
	err := util.LoadConfig(".")
	if err != nil {
		return
	}
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, util.Config.DBSource)
	if err != nil {
		panic(err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(util.Config.Port)
	if err != nil {
		return
	}
	defer conn.Close()
}