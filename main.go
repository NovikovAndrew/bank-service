package main

import (
	"bank-service/token"
	"bank-service/util"
	"database/sql"
	"log"

	"bank-service/api"
	"bank-service/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	pasetoMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(store, pasetoMaker, config)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal(err)
	}
}
