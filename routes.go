package main

import (
	"net/http"

	"recordstore-go/gen/recordstore/v1/recordstorev1connect"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	connectPath, connectHandler := recordstorev1connect.NewRecordStoreServiceHandler(app)
	mux.Handle(connectPath, connectHandler)

	// OAuth callbacks stay as plain HTTP — providers redirect browsers here.
	mux.HandleFunc("/v1/spotify/callback", app.spotifyCallbackHandler)
	mux.HandleFunc("/v1/amazon/callback", app.amazonCallbackHandler)

	return app.enableCORS(mux)
}
