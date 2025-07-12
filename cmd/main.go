package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/db/clients"
	"github.com/bercivarga/go-basic-server/internal/middleware"
	"github.com/bercivarga/go-basic-server/internal/router"
	"github.com/bercivarga/go-basic-server/internal/wire"
)

const (
	defaultDSN  = "localSQLite.db"
	defaultPort = 8080
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := flag.Int("port", defaultPort, "Port to run the server on")
	flag.Parse()

	sqlite := clients.NewSQLite(defaultDSN)

	if _, err := sqlite.Connect(); err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer func(sqlite *clients.SQLite) {
		err := sqlite.Close()
		if err != nil {
			log.Fatalf("DB close: %v", err)
		}
	}(sqlite)

	app := app.NewApp(sqlite.DB)

	router := router.New(app)

	wire := wire.New(app)
	wire.RegisterRoutes(router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      middleware.Logger(router),
	}

	log.Printf("Starting server on port %d", *port)
	serverStartErr := server.ListenAndServe()
	if serverStartErr != nil {
		log.Fatalf("Server listen: %v", serverStartErr)
	}
}
