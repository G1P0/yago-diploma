package main

import (
	"log"
	"os"

	"github.com/G1P0/yago-diploma/internal/server"
	"github.com/G1P0/yago-diploma/pkg/db"
)

func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags)

	if err := db.Init("scheduler.db"); err != nil {
		logger.Fatalf("failed to init db: %v", err)
	}

	srv := server.NewServer(logger)

	if err := srv.Server.ListenAndServe(); err != nil {
		srv.Logger.Fatalf("error while starting server: %s", err)
	}

}
