package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recordstore-go/adapters"
	"recordstore-go/models"
	"time"
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
	fmt.Println("spotifyUserMusicDataHandler Get sptfySession")
	sptfySession := r.URL.Query().Get("sptfySession")

	// Check if a current session
	currentSession, _ := models.IsCurrentSession(cfg.db.conn, sptfySession)

	userMusicData := models.MusicData{}
	fmt.Println("spotifyUserMusicDataHandler Check if new")
	if !currentSession && sptfySession != "" {
		// Check if user exists in system
		// This is not a current session 
		// spotifyID, err := models.GetSpotifyIDBySessionID(sptfySession)
		// if err != nil {
		// 	fmt.Println("GetSpotifyIDBySessionID error")
		// 	return
		// }
		// user, err := models.GetUserBySpotifyID(cfg.db.conn, spotifyID)
		// if err != nil {
		// 	fmt.Println("GetUserBySpotifyID error")
		// 	return
		// }

		// if user != nil {
		// 	// Return current stored data

		// 	followedArtists, err := user.GetUserArtists(cfg.db.conn)
		// 	savedAlbums, err := user.GetUserArtists(cfg.db.conn)

		// 	fmt.Println("spotifyUserMusicDataHandler Bundle Data")
		// 	userMusicData = models.MusicData{
		// 		Albums:   savedAlbums,
		// 		Artists:  followedArtists,
		// 		Analysis: user.Analysis,
		// 	}

		// } else 
		if true {
			fmt.Println("spotifyUserMusicDataHandler Get token")
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
			// if cfg.env != "develop" {
			if true {
				vertexAI := adapters.NewAdapter("")
				vertexParams := map[string]interface{}{
					"temperature":     0.2,
					"maxOutputTokens": 500,
					"topP":            0.95,
					"topK":            40,
				}
				err, musicAnalysis = vertexAI.TextPredict(w, followedArtists, "hipster-record-store-clerk", cfg.vertex, vertexParams)
				if err != nil {
					fmt.Println("ArtistsResponseRequest error", err)
				}
			} else {
				musicAnalysis = "No textPredict data in develop mode"
			}

			fmt.Println("spotifyUserMusicDataHandler Bundle Data")
			userMusicData = models.MusicData{
				Albums:   savedAlbums,
				Artists:  followedArtists,
				Analysis: musicAnalysis,
			}

			// // Need to spin this off to update albums and artists and relationships
			// // go
			// userMusicData.Artists.InsertAll(cfg.db.conn)

			// Get the spotify user info
			adapter.GetSpotifyUserData(sptfySession)

			// Create Spotify Session
			user := models.User{}
			err = user.CreateUserSpotifySession(cfg.db.conn, sptfySession)
			if err != nil {
				fmt.Println("spotifyUserMusicDataHandler:CreateUserSpotifySession error")
				return
			}

			// Need to spin this off to update albums and artists and relationships
			go user.UpdateUserLibrary(cfg.db.conn, userMusicData)
			// user.AddUserArtists(cfg.db.conn, followedArtists)
		}

	} else {
		fmt.Println("spotifyUserMusicDataHandler Existing")
		// Get existing data

		// userMusicData = models.MusicData{
		// 	Albums:   savedAlbums,
		// 	Artists:  followedArtists,
		// 	Analysis: musicAnalysis,
		// }
	}

	js, _ := json.MarshalIndent(userMusicData, "", "\t")

	// JV TODO: Temp sessData
	fmt.Println("spotifyUserMusicDataHandler Create cookie")
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:    "userSession",
		Value:   "1",
		Expires: expiration}

	fmt.Println("spotifyUserMusicDataHandler Set cookie")
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}
