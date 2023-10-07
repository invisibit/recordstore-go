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
	err, sptfyToken := adapter.GetSpotifyUserAccessToken(code, cfg.client_id, cfg.client_secret)
	if err != nil {
		fmt.Println("GetSpotifyUserAccessToken error")
		return
	}

	// js, err := json.MarshalIndent(followedArtists, "", "\t")

	// sessionCookie := "sptfySession=" + sptfyToken
	http.Redirect(w, r, "http://localhost:3000/Mymusic?sptfySession=" + sptfyToken, http.StatusSeeOther)

}
