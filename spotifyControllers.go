package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recordstore-go/adapters"
)

// TODO: These functions are deprecated and can be removed.

func (app *application) spotifyFollowedArtistsHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("spotifyFollowedArtistsHandler")
	// Get parameters
	sptfySession := r.URL.Query().Get("sptfySession")

	// Use the token to get followed artists
	adapter := adapters.NewAdapter("")
	err, followedArtists := adapter.GetSpotifyUserFollowedArtists(sptfySession)
	if err != nil {
		fmt.Println("GetSpotifyUserFollowedArtists error")
		return
	}

	js, err := json.MarshalIndent(followedArtists, "", "\t")

	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func (app *application) spotifySavedAlbumsHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("spotifySavedAlbumsHandler")
	// Get parameters
	sptfySession := r.URL.Query().Get("sptfySession")

	// Use the token to get saved Albums
	adapter := adapters.NewAdapter("")
	err, savedAlbums := adapter.GetSpotifyUserSavedAlbums(sptfySession)
	if err != nil {
		fmt.Println("spotifySavedAlbumsHandler error")
		return
	}

	js, err := json.MarshalIndent(savedAlbums, "", "\t")

	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func (app *application) spotifyDevicePlayback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter spotifyDevicePlayback")
	albumID := r.URL.Query().Get("albumID")
	deviceID := r.URL.Query().Get("audioDeviceID")
	sptfySession := r.URL.Query().Get("sptfySession")

	adapter := adapters.NewAdapter("")
	err := adapter.PlaySpotifyAlbum(sptfySession, albumID, deviceID)
	if err != nil {
		fmt.Println("PlaySpotifyAlbum error")
		return
	}

}
