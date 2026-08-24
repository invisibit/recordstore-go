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
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


// TODOS - 20240103
// Add route for existing user
// Move controllers to controller directory
// Add text in UI to update a service data click on the service button

const version = "1.0.1"

type Config struct {
	srv_addr              string
	srv_port              string
	env                   string
	client_id             string
	client_secret         string
	youtube_client_id     string
	youtube_client_secret string
	ui_address            string
	ui_port               string

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
	if err := godotenv.Load(); err != nil && cfg.env != "prod" {
		log.Println("No .env file found, relying on real environment variables")
	}

	cfg.env = os.Getenv("env")
	cfg.srv_addr = os.Getenv("srv_address")
	cfg.srv_port = os.Getenv("srv_port")
	cfg.client_id = os.Getenv("client_id")
	cfg.client_secret = os.Getenv("client_secret")
	cfg.youtube_client_id = os.Getenv("youtube_client_id")
	cfg.youtube_client_secret = os.Getenv("youtube_client_secret")
	cfg.ui_address = os.Getenv("ui_address")
	cfg.ui_port = os.Getenv("ui_port")
	cfg.db.dsn = os.Getenv("dsn")
	cfg.db.user = os.Getenv("db_user")
	cfg.db.password = os.Getenv("db_password")
	cfg.db.dbName = os.Getenv("db_name")
	cfg.db.dbPort = os.Getenv("db_port")
	cfg.vertex.Location = os.Getenv("vertexai_location")
	cfg.vertex.Publisher = os.Getenv("vertexai_publisher")
	cfg.vertex.Model = os.Getenv("vertexai_model")

	flag.StringVar(&cfg.db.dsn, "dsn", "", "Postgres connection string (overrides host/user/password/dbname)")
	flag.Parse()

	// In Cloud Run, PORT is set to 8080. Default srv_port if unset.
	if cfg.srv_port == "" {
		cfg.srv_port = os.Getenv("PORT")
	}
	if cfg.srv_port == "" {
		cfg.srv_port = "8080"
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	// Set up database connection
	dsn := cfg.db.dsn
	if dsn == "" {
		dsn = "host=" + os.Getenv("db_host") +
			" user=" + cfg.db.user +
			" password=" + cfg.db.password +
			" dbname=" + cfg.db.dbName +
			" port=" + cfg.db.dbPort +
			" sslmode=disable"
	}
	if dsn != "" {
		conn, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		// defer conn.Close()
		cfg.db.conn = conn

		if err != nil {
			log.Println(err)
		}

		// Migrate the schema(probably move to a seperate function)
		_ = conn.Exec("CREATE DATABASE IF NOT EXISTS record_store_clerk;")
		conn.AutoMigrate(&models.Album{})
		conn.AutoMigrate(&models.Artist{})
		conn.AutoMigrate(&models.User{})
		conn.AutoMigrate(&models.ConnectedProvider{})
		conn.AutoMigrate(&models.UserArtist{})
		conn.AutoMigrate(&models.UserAlbum{})

	}

	// hostname := ""
	// if cfg.env != "develop" {
	// 	hostname = "localhost"
	// } else {
	// 	hostname = "0.0.0.0"
	// }
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.srv_addr, cfg.srv_port),
		Handler:      h2c.NewHandler(app.routes(), &http2.Server{}),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server open port", cfg.srv_port)
	fmt.Println("Starting server at", cfg.srv_addr, cfg.srv_port)

	srvErr := srv.ListenAndServe()
	if srvErr != nil {
		log.Println(srvErr)
	}
}
