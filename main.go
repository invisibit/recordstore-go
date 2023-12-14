package main

import (
	// "recordstore-go/models"
	"recordstore-go/adapters"

	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const version = "1.0.1"

type config struct {
	port          int
	env           string
	client_id     string
	client_secret string
	ui_address    string

	db struct {
		dsn      string
		user     string
		password string
		dbName   string
		dbPort   string
	}

	vertex			adapters.VertexModelParams
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

	if cfg.env != "develop" {
		fmt.Println("Load .env")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file AAA")
		}
	}
	
	cfg.client_id = os.Getenv("client_id")
	cfg.client_secret = os.Getenv("client_secret")
	cfg.ui_address = os.Getenv("ui_address")
	cfg.env = os.Getenv("env")
	cfg.db.dsn = os.Getenv("dsn")
	cfg.db.user = os.Getenv("db_user")
	cfg.db.password = os.Getenv("db_password")
	cfg.db.dbName = os.Getenv("db_name")
	cfg.vertex.Location = os.Getenv("vertexai_location")
	cfg.vertex.Publisher = os.Getenv("vertexai_publisher")
	cfg.vertex.Model = os.Getenv("vertexai_model")

	// flag.IntVar(&cfg.port, "port", 4000, "server port listen on ")
	// flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	// flag.StringVar(&cfg.db.dsn, "dsn", "postgres://root:root@127.0.0.1:5434/testingwithrentals?sslmode=disable", "Postgres connection string")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	// Set up database connection
	if cfg.db.dsn != "" {
		_, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=localhost user=" + cfg.db.user + " password=" + cfg.db.user + " dbname=" + cfg.db.dbName + " port=" + cfg.db.dbPort + " sslmode=disable", // data source name, refer https://github.com/jackc/pgx
			PreferSimpleProtocol: true,                                                                                                                                           // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
		}), &gorm.Config{})
		// defer db.close()

		if err != nil {
			log.Println(err)
		}
	}

	// Migrate the schema(probably move to a seperate function)
	// db.AutoMigrate(&models.Artist{})

	hostname := ""
	if cfg.env != "develop" {
		hostname = "localhost"
	} else {
		hostname = "0.0.0.0"
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", hostname, cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server open port", cfg.port)

	srvErr := srv.ListenAndServe()
	if srvErr != nil {
		log.Println(srvErr)
	}
}
