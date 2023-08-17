package main

import (
	// "backend/cmd/api/models"

	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.1"

type config struct {
	port          int
	env           string
	client_id     string
	client_secret string

	// db   struct {
	// 	dsn string
	// }
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	// models models.Models
}

var cfg config

func main() {

	// Start up the server
	fmt.Println("Starting server")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.client_id = os.Getenv("client_id")
	cfg.client_secret = os.Getenv("client_secret")

	flag.IntVar(&cfg.port, "port", 4000, "server port listen on ")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server open port", cfg.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
