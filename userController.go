package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recordstore-go/models"
	"strconv"
)

func (app *application) userHandler(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("userSession")
	fmt.Println("Cookie", cookie)

	// Get user information
	userID, _ := strconv.ParseUint(cookie.Value, 10, 64)
	user := models.User{ID: uint(userID)}

	js, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}

	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func (app *application) userLibraryHandler(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("userSession")
	fmt.Println("Cookie", cookie)
	userID, _ := strconv.ParseUint(cookie.Value, 10, 64)
	user := models.User{ID: uint(userID)}

	// Get user library by getting the users albums then the artists
	albums, err := user.GetUserAlbums(cfg.db.conn)
	if err != nil {
		app.logger.Println(err)
	}

	// Get user library by getting the users albums then the artists
	artists, err := user.GetUserArtists(cfg.db.conn)
	if err != nil {
		app.logger.Println(err)
	}

	userMusicData := models.MusicData{}
	userMusicData.Albums = albums
	userMusicData.Artists = artists

	js, err := json.MarshalIndent(userMusicData, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}

	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}
