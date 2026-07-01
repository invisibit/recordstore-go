package main

import (
	"fmt"
	"net/http"
	"recordstore-go/adapters"
	"time"
)

func (app *application) spotifyCallbackHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("spotifyCallbackHandler")
	// Get parameters
	code := r.URL.Query().Get("code")
	loginError := r.URL.Query().Get("error")
	// state := r.URL.Query().Get("state")

	fmt.Println("code", code)
	if loginError != "" {
		fmt.Println("loginError", loginError)
	}

	// Use the code to get the access token
	adapter := adapters.NewAdapter("https://accounts.spotify.com/")
	redirectHost := ""
	if cfg.env == "develop" {
		redirectHost = "http://localhost:4000/v1/spotify/callback"
	} else {
		redirectHost = "https://" + r.Host + "/v1/spotify/callback"
	}

	err, sptfyToken := adapter.GetSpotifyUserAccessToken(code, cfg.client_id, cfg.client_secret, redirectHost)
	if err != nil {
		fmt.Println("GetSpotifyUserAccessToken error")
		return
	}

	fmt.Println("--------------------------------------------Redirect---------------------")

	if cfg.env == "develop" {
		fmt.Println("spotifyCallbackHandler Create cookie")
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:    "userSession",
			Value:   "1",
			Expires: expiration}

		fmt.Println("spotifyCallbackHandler Set cookie")
		http.SetCookie(w, &cookie)
		fmt.Println("Develop environment")
		http.Redirect(w, r, "http://localhost:3000/Mymusic?sptfySession="+sptfyToken, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "https://"+cfg.ui_address+"/Mymusic?sptfySession="+sptfyToken, http.StatusSeeOther)
	}

}

func (app *application) amazonCallbackHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("amazonCallbackHandler")
	// Get parameters
	code := r.URL.Query().Get("code")
	loginError := r.URL.Query().Get("error")

	fmt.Println("code", code)
	if loginError != "" {
		fmt.Println("loginError", loginError)
	}

}
