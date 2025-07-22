package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error

	testDB, err = pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	testQueries = New(testDB)

	defer testDB.Close()

	os.Exit(m.Run())
}
