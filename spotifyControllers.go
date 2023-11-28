package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recordstore-go/adapters"
	"recordstore-go/models"
)

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

func (app *application) spotifyUserMusicDataHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("spotifyUserMusicDataHandler")
	// Get parameters
	sptfySession := r.URL.Query().Get("sptfySession")

	// Use the token to get followed artists
	adapter := adapters.NewAdapter("")
	err, followedArtists := adapter.GetSpotifyUserFollowedArtists(sptfySession)
	if err != nil {
		fmt.Println("GetSpotifyUserFollowedArtists error")
		return
	}

	// Get the saved albums
	err, savedAlbums := adapter.GetSpotifyUserSavedAlbums(sptfySession)
	if err != nil {
		fmt.Println("spotifySavedAlbumsHandler error")
		return
	}

	// Get the analysis of your artists
	musicAnalysis := ""
	if true {
		vertexAI := adapters.NewAdapter("")
		vertexParams := map[string]interface{}{
			"temperature":     0.2,
			"maxOutputTokens": 500,
			"topP":            0.95,
			"topK":            40,
		}
		err, musicAnalysis = vertexAI.TextPredict(w, followedArtists, "hipster-record-store-clerk", "us-central1", "google", "text-bison@001", vertexParams)
		if err != nil {
			fmt.Println("ArtistsResponseRequest error", err)
			return
		}
	}

	userMusicData := models.MusicData{
		Albums:   savedAlbums,
		Artists:  followedArtists,
		Analysis: musicAnalysis,
	}

	js, err := json.MarshalIndent(userMusicData, "", "\t")

	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}
