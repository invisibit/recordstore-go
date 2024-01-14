package main

import (
	"recordstore-go/adapters"
	"recordstore-go/models"

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

// TODOS - 20240103
// Add route for existing user
// Move controllers to controller directory
// Add text in UI to update a service data click on the service button

const version = "1.0.1"

type Config struct {
	port          string
	env           string
	client_id     string
	client_secret string
	ui_address    string

	db struct {
		conn     *gorm.DB
		dsn      string
		user     string
		password string
		dbName   string
		dbPort   string
	}

	vertex adapters.VertexModelParams
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config Config
	logger *log.Logger
	// models models.Models
}

var cfg Config

func main() {

	// Start up the server
	fmt.Println("Starting server")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	fmt.Println("Load .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file Add .env to prod that's blank")
	}

	cfg.env = os.Getenv("env")
	cfg.port = os.Getenv("srv_port")
	cfg.client_id = os.Getenv("client_id")
	cfg.client_secret = os.Getenv("client_secret")
	cfg.ui_address = os.Getenv("ui_address")
	cfg.db.dsn = os.Getenv("dsn")
	cfg.db.user = os.Getenv("db_user")
	cfg.db.password = os.Getenv("db_password")
	cfg.db.dbName = os.Getenv("db_name")
	cfg.db.dbPort = os.Getenv("db_port")
	cfg.vertex.Location = os.Getenv("vertexai_location")
	cfg.vertex.Publisher = os.Getenv("vertexai_publisher")
	cfg.vertex.Model = os.Getenv("vertexai_model")

	// flag.IntVar(&cfg.port, "port", 4000, "server port listen on ")
	// flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	// flag.StringVar(&cfg.db.dsn, "dsn", "postgres://root:root@127.0.0.1:5434/testingwithrentals?sslmode=disable", "Postgres connection string")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://root:root@127.0.0.1:5434/testingwithrentals?sslmode=disable", "Postgres connection string")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	// Set up database connection
	if cfg.db.dsn != "" {
		conn, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=localhost user=" + cfg.db.user + " password=" + cfg.db.password + " dbname=" + cfg.db.dbName + " port=" + cfg.db.dbPort + " sslmode=disable", // data source name, refer https://github.com/jackc/pgx
			PreferSimpleProtocol: true,                                                                                                                                               // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
		}), &gorm.Config{})
		// defer conn.Close()
		cfg.db.conn = conn

		if err != nil {
			log.Println(err)
		}

		// Migrate the schema(probably move to a seperate function)
		_ = conn.Exec("CREATE DATABASE IF NOT EXISTS record_store;")
		conn.AutoMigrate(&models.Album{})
		conn.AutoMigrate(&models.Artist{})
		conn.AutoMigrate(&models.User{})
		conn.AutoMigrate(&models.UserArtist{})
		conn.AutoMigrate(&models.UserAlbum{})

	}

	hostname := ""
	if cfg.env != "develop" {
		hostname = "localhost"
	} else {
		hostname = "0.0.0.0"
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", hostname, cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server open port", cfg.port)

	// fmt.Println("Config struct", cfg)

	srvErr := srv.ListenAndServe()
	if srvErr != nil {
		log.Println(srvErr)
	}
}
