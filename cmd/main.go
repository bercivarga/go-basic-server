package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/db/clients"
	"github.com/bercivarga/go-basic-server/internal/logger"
	"github.com/bercivarga/go-basic-server/internal/router"
	"github.com/bercivarga/go-basic-server/internal/wire"
)

const (
	defaultDSN  = "localSQLite.db"
	defaultPort = 8080
)

func main() {
	port := flag.Int("port", defaultPort, "Port to run the server on")
	flag.Parse()

	sqlite := clients.NewSQLite(defaultDSN)

	if _, err := sqlite.Connect(); err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer sqlite.Close()

	app := app.NewApp(sqlite.DB, logger.New())

	router := router.New(app)

	wire := wire.New(app)
	wire.RegisterRoutes(router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      logger.Middleware(router),
	}

	log.Printf("Starting server on port %d", *port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server listen: %v", err)
	}
}
