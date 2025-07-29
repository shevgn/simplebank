// Package main provides the main entry point for the simplebank application.
package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shevgn/simplebank/api"
	db "github.com/shevgn/simplebank/db/sqlc"
	"github.com/shevgn/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(config, store)

	err = server.Run(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
