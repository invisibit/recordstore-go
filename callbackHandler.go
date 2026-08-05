package main

import (
	"fmt"
	"net/http"
	"os"
	"recordstore-go/adapters"
	"recordstore-go/models"
	"strconv"
	"strings"
	"time"
)

func setUserSessionCookie(w http.ResponseWriter, userID uint) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "userSession",
		Value:    strconv.FormatUint(uint64(userID), 10),
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func GetUserSessionCookie(r *http.Request) uint64 {
	cookie, err := r.Cookie("userSession")
	if err != nil {
		return 0
	}

	id, err := strconv.ParseUint(cookie.Value, 10, 64)
	if err != nil {
		return 0
	}

	return id
}

// Build the UI address
// Changed from devUIAddress to getUIAddress to handle both dev and prod environments
func getUIAddress(callbackHost string) string {
	port := os.Getenv("ui_port")
	host := callbackHost
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return "https://" + host + ":" + port
}

func (app *application) spotifyCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("spotifyCallbackHandler")
	userID := GetUserSessionCookie(r)
	user := models.User{}
	fmt.Println("spotifyCallbackHandler", userID)
	if userID != 0 {
		// Get existing user by ID if userID is provided
		existingUser, err := models.GetExistingSessionByUserID(cfg.db.conn, strconv.Itoa(int(userID)))
		if err != nil {
			fmt.Println("spotifyCallbackHandler: GetUserByID error", err)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		user = existingUser
	} else {
		// Create a new user if no userID is provided
		user.CreateUser(cfg.db.conn)
	}
	// fmt.Println("spotifyCallbackHandler", user.ID)

	code := r.URL.Query().Get("code")
	loginError := r.URL.Query().Get("error")

	// fmt.Println("spotifyCallbackHandler code", code)
	if loginError != "" {
		fmt.Println("loginError", loginError)
	}

	adapter := adapters.NewAdapter("https://accounts.spotify.com/")
	redirectHost := ""
	if cfg.env == "develop" {
		redirectHost = "https://127.0.0.1:4000/v1/spotify/callback"
	} else {
		redirectHost = "https://" + r.Host + "/v1/spotify/callback"
	}

	err, sptfyToken := adapter.GetSpotifyUserAccessToken(code, cfg.client_id, cfg.client_secret, redirectHost)
	if err != nil {
		fmt.Println("GetSpotifyUserAccessToken error")
		http.Error(w, "spotify token exchange failed", http.StatusInternalServerError)
		return
	}

	// If usernam is blank use the spotify username
	if user.UserName == "" {
		user.UserName, err = adapter.GetSpotifyUserDisplayName(user.SpotifySession)
		if err != nil {
			fmt.Println("spotifyCallbackHandler GetSpotifyUserDisplayName error", err)
		}
	}
	// Update with new spotify token
	user.SpotifySession = sptfyToken

	if cfg.db.conn != nil {
		err = user.UpdateUser(cfg.db.conn)
		if err != nil {
			fmt.Println("spotifyCallbackHandler UpdateUser error", err)
		} else {
			setUserSessionCookie(w, user.ID)
		}
	}

	uiAddress := getUIAddress(r.Host)
	http.Redirect(w, r, uiAddress+"/minimax", http.StatusSeeOther)
}

func (app *application) youtubeCallbackHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetUserSessionCookie(r)
	user := models.User{}
	fmt.Println("youtubeCallbackHandler", userID)
	if userID != 0 {
		// Get existing user by ID if userID is provided
		existingUser, err := models.GetExistingSessionByUserID(cfg.db.conn, strconv.Itoa(int(userID)))
		if err != nil {
			fmt.Println("spotifyCallbackHandler: GetUserByID error", err)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		user = existingUser
	} else {
		// Create a new user if no userID is provided
		user.CreateUser(cfg.db.conn)
	}

	// fmt.Println("youtubeCallbackHandler", r.URL.Query().Get("code"))
	code := r.URL.Query().Get("code")
	loginError := r.URL.Query().Get("error")

	// fmt.Println("youtubeCallbackHandler code", code)
	if loginError != "" {
		fmt.Println("loginError", loginError)
		http.Error(w, "google login error: "+loginError, http.StatusBadRequest)
		return
	}
	if code == "" {
		fmt.Println("missing code in youtube callback")
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}
	adapter := adapters.NewAdapter("https://oauth2.googleapis.com/token")

	redirectHost := ""
	if cfg.env == "develop" {
		redirectHost = "https://127.0.0.1:4000/v1/youtube/callback"
	} else {
		redirectHost = "https://" + r.Host + "/v1/youtube/callback"
	}

	err, youtubeToken := adapter.GetYoutubeUserAccessToken(code, cfg.youtube_client_id, cfg.youtube_client_secret, redirectHost)
	if err != nil {
		fmt.Println("GetYoutubeUserAccessToken error", err)
		http.Error(w, "youtube token exchange failed", http.StatusInternalServerError)
		return
	}

	// fmt.Println("youtube access_token", youtubeToken.Access_token)
	// fmt.Println("youtube refresh_token", youtubeToken.Refresh_token)

	// Get the user name from Google if there isn't one currently set
	if user.UserName == "" {
		user.UserName, err = adapter.GetGoogleUserDisplayName(youtubeToken.Access_token)
		if err != nil {
			fmt.Println("GetGoogleUserDisplayName error", err)
		}
	}

	// Update the user with the new YouTube session and refresh tokens
	user.YoutubeSession = youtubeToken.Access_token
	user.YoutubeRefreshToken = youtubeToken.Refresh_token

	if cfg.db.conn != nil {
		err = user.UpdateUser(cfg.db.conn)
		if err != nil {
			fmt.Println("GetOrCreateUserByUserID(youtube) error", err)
		} else {
			setUserSessionCookie(w, user.ID)
		}
	}

	uiAddress := getUIAddress(r.Host)
	http.Redirect(w, r, uiAddress+"/minimax?provider=youtube", http.StatusSeeOther)
}

func (app *application) amazonCallbackHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("amazonCallbackHandler")
	code := r.URL.Query().Get("code")
	loginError := r.URL.Query().Get("error")

	fmt.Println("code", code)
	if loginError != "" {
		fmt.Println("loginError", loginError)
	}

}
