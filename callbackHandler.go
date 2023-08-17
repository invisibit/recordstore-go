package main

import (
	"fmt"
	"net/http"
	"recordstore-go/adapters"
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
	adapter.GetSpotifyUserAccessToken(code, cfg.client_id, cfg.client_secret)

	// adapter.GetSpotifyUserData(code)
	// fmt.Println("state", state)

	w.WriteHeader(http.StatusOK)
	// w.Write("Thanks")
}
