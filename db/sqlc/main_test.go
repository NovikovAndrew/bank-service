package db

import (
	"bank-service/util"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQuery *Queries
	testDB    *sql.DB
)

func TestMain(m *testing.M) {
	config, configErr := util.LoadConfig("../../")
	if configErr != nil {
		log.Fatal(configErr)
	}

	var err error
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQuery = New(testDB)

	os.Exit(m.Run())
}
