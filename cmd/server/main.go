package main

import (
	"context"
	"log"
	"os"

	"github.com/example/ledger/internal/api"
	"github.com/example/ledger/internal/db"
	"github.com/example/ledger/internal/repository"
	"github.com/example/ledger/internal/service"
)

func main() {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://root:secret@localhost:5432/ledger?sslmode=disable"
	}

	connPool, err := db.NewConnectionPool(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer connPool.Close()

	store := repository.NewStore(connPool)
	svc := service.NewService(store)
	server := api.NewServer(svc)

	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
