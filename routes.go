package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/spotify/callback", app.spotifyCallbackHandler)
	router.HandlerFunc(http.MethodGet, "/v1/spotify/followed", app.spotifyFollowedArtistsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/spotify/savedAlbums", app.spotifySavedAlbumsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/spotify/userMusicData", app.spotifyUserMusicDataHandler)

	return app.enableCORS(router)
}
