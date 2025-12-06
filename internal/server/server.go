package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/G1P0/yago-diploma/internal/handlers"
	"github.com/G1P0/yago-diploma/pkg/api"
)

type Server struct {
	Logger *log.Logger
	Server *http.Server
}

func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HandleRoot)

	api.Init(mux)

	addr := ":7540"

	if port := os.Getenv("TODO_PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: logger,
		Server: srv,
	}
}
