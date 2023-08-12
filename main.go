package main

import (
	// "backend/cmd/api/models"

	"recordstore-go/adapters"

	"flag"
	"fmt"
	"net/http"
	"log"
	"os"
	"time"
	_ "github.com/lib/pq"
)

const version = "1.0.1"

type config struct {
	port int
	env  string
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

func main() {

	// Start up the server
	fmt.Println("Starting server")

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "server port listen on ")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server open port", cfg.port)

	logger.Println("Spotify Login")
	adapter := adapters.NewAdapter("https://accounts.spotify.com/api/token")

	err := adapter.OpenSpotifyConnection()
	if err != nil {
		log.Println(err)
	}

	err = adapter.OpenSpotifyUserConnection()
	if err != nil {
		log.Println(err)
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	// Get user info

}

// // OpenSpotifyConnection gets a bearer to use for this session
// func OpenSpotifyConnection() error {
// 		// Create the endpoint from the task
// 		// endPoint := t.Type + "/" + t.ID
//         // fmt.Printf("worker  %d retrieving %s\n", id, endPoint)

// 		// Retrieve token from api
// 		urlReguest := "https://accounts.spotify.com/api/token"

// 		parm := url.Values{}
// 		parm.Add("grant_type", "client_credentials")
// 		parm.Add("client_id", "5c37b6f4c90143908b11d9e1727db5e7")
// 		parm.Add("client_secret", "0a5f8bc725a94eb7aae28803adc108d1")

// 	   cookieJar, _ := cookiejar.New(nil)
// 		tr := &http.Transport{
// 			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 		}
// 		client := &http.Client{Transport: tr,
// 			Jar: cookieJar,
// 			CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 				return http.ErrUseLastResponse
// 			}}
		
// 		req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
// 		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Printf("%s", err)
// 		}
// 		// fmt.Println(resp.Header)
// 		defer resp.Body.Close()
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Printf("%s", err)
// 			return err
// 		}
// 		fmt.Println(string(body))

// 		return nil

// 		// resp, err := http.Get(urlReguest)
// 		// if err != nil {1
// 		// 	log.Fatal("worker error ", urlReguest)
// 		// 	return
// 		// }

// }
