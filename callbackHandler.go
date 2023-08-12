package main

import (
	"fmt"
	"net/http"
)

func (app *application) spotifyCallbackHandler(w http.ResponseWriter, r *http.Request) {

	// Get parameters
	code := r.URL.Query().Get("code")
	loginError := r.URL.Query().Get("error")
	state := r.URL.Query().Get("state")

	fmt.Println("code", code)
	fmt.Println("loginError", loginError)
	fmt.Println("state", state)
}
