package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"connectrpc.com/connect"

	"recordstore-go/adapters"
	v1 "recordstore-go/gen/recordstore/v1"
	"recordstore-go/models"
)

const userIDHeader = "X-User-Id"
const userSessionCookie = "userSession"

func userIDFromCookieHeader(req connect.AnyRequest) uint64 {
	cookieHeader := req.Header().Get("Cookie")
	if cookieHeader == "" {
		return 0
	}

	cookies, err := http.ParseCookie(cookieHeader)
	if err != nil {
		return 0
	}

	for _, cookie := range cookies {
		if cookie.Name != userSessionCookie {
			continue
		}
		id, err := strconv.ParseUint(cookie.Value, 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func currentUserID(req connect.AnyRequest) uint64 {
	cookieID := userIDFromCookieHeader(req)
	if cookieID != 0 {
		return cookieID
	}
	v := req.Header().Get(userIDHeader)
	if v == "" {
		return 0
	}
	id, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func resolveProviderTokens(userID uint64) (string, string) {
	if userID == 0 || cfg.db.conn == nil {
		return "", ""
	}
	user, err := models.GetUserByID(cfg.db.conn, uint(userID))
	if err != nil {
		return "", ""
	}
	return user.SpotifySession, user.YoutubeSession
}

func (app *application) GetCurrentUser(
	ctx context.Context,
	req *connect.Request[v1.GetCurrentUserRequest],
) (*connect.Response[v1.GetCurrentUserResponse], error) {
	// Diagnostic logging to help debug cookie/header presence
	// fmt.Println("GetCurrentUser request headers:", req.Header())
	id := currentUserID(req)
	fmt.Println("Enter GetCurrentUser id", id)
	if id == 0 || cfg.db.conn == nil {
		return connect.NewResponse(&v1.GetCurrentUserResponse{}), nil
	}

	user, err := models.GetUserByID(cfg.db.conn, uint(id))
	if err != nil {
		return connect.NewResponse(&v1.GetCurrentUserResponse{}), nil
	}

	cps, err := user.GetConnectedProviders(cfg.db.conn)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	providers := make([]string, 0, len(cps))
	for _, cp := range cps {
		providers = append(providers, cp.Provider)
	}

	return connect.NewResponse(&v1.GetCurrentUserResponse{
		UserId:    id,
		Username:  user.UserName,
		Providers: providers,
	}), nil
}

func (app *application) Status(
	ctx context.Context,
	req *connect.Request[v1.StatusRequest],
) (*connect.Response[v1.StatusResponse], error) {
	return connect.NewResponse(&v1.StatusResponse{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}), nil
}

func (app *application) GetFollowedArtists(
	ctx context.Context,
	req *connect.Request[v1.GetFollowedArtistsRequest],
) (*connect.Response[v1.GetFollowedArtistsResponse], error) {
	adapter := adapters.NewAdapter("")
	err, artists := adapter.GetSpotifyUserFollowedArtists(req.Msg.SptfySession)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.GetFollowedArtistsResponse{
		Artists: toProtoArtists(artists),
	}), nil
}

// GetSavedAlbums retrieves the user's saved albums from all providers.
func (app *application) GetSavedAlbums(ctx context.Context, req *connect.Request[v1.GetSavedAlbumsRequest]) (*connect.Response[v1.GetSavedAlbumsResponse], error) {
	// user := models.User{}
	user, err := models.GetUserByID(cfg.db.conn, uint(currentUserID(req)))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	adapter := adapters.NewAdapter("")
	albums := []models.Album{}
	if user.SpotifySession != "" {

		err, spotifyAlbums := adapter.GetSpotifyUserSavedAlbums(user.SpotifySession)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		albums = append(albums, spotifyAlbums...)
	}

	if user.YoutubeSession != "" {
		err, youtubeAlbums := adapter.GetYoutubeUserSavedAlbums(user.YoutubeSession)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		albums = append(albums, youtubeAlbums...)
	}

	return connect.NewResponse(&v1.GetSavedAlbumsResponse{
		Albums: toProtoAlbums(albums),
	}), nil
}

func (app *application) GetYoutubeSavedAlbums(ctx context.Context, req *connect.Request[v1.GetYoutubeSavedAlbumsRequest]) (*connect.Response[v1.GetYoutubeSavedAlbumsResponse], error) {
	youtubeSession := req.Msg.YoutubeSession
	if youtubeSession == "" {
		youtubeSession = req.Header().Get("youtubeSession")
	}
	if youtubeSession == "" {
		_, youtubeSession = resolveProviderTokens(currentUserID(req))
	}
	if youtubeSession == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no youtube session available"))
	}

	adapter := adapters.NewAdapter("")
	err, albums := adapter.GetYoutubeUserSavedAlbums(youtubeSession)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.GetYoutubeSavedAlbumsResponse{
		Albums: toProtoAlbums(albums),
	}), nil
}

func (app *application) GetUserMusicData(
	ctx context.Context,
	req *connect.Request[v1.GetUserMusicDataRequest],
) (*connect.Response[v1.GetUserMusicDataResponse], error) {
	fmt.Println("Enter GetMusicUserData")
	adapter := adapters.NewAdapter("")

	err, artists := adapter.GetSpotifyUserFollowedArtists(req.Msg.SptfySession)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// TODO add after streaming worked out
	// err, albums := adapter.GetSpotifyUserSavedAlbums(req.Msg.SptfySession)
	// if err != nil {
	// 	return nil, connect.NewError(connect.CodeInternal, err)
	// }
	albums := []models.Album{}

	analysis := "No textPredict data in develop mode"
	if app.config.env != "develop" {
		vertexAI := adapters.NewAdapter("")
		vertexParams := map[string]interface{}{
			"temperature":     0.2,
			"maxOutputTokens": 500,
			"topP":            0.95,
			"topK":            40,
		}
		err, analysis = vertexAI.TextPredict(io.Discard, artists, "hipster-record-store-clerk", app.config.vertex, vertexParams)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&v1.GetUserMusicDataResponse{
		Artists:  toProtoArtists(artists),
		Albums:   toProtoAlbums(albums),
		Analysis: analysis,
	}), nil
}

func toProtoArtists(in []models.Artist) []*v1.Artist {
	out := make([]*v1.Artist, len(in))
	for i, a := range in {
		out[i] = &v1.Artist{
			Id:            a.SpotifyID,
			SpotifyId:     a.SpotifyID,
			Name:          a.Name,
			ExternalUrls:  a.ExternalUrls,
			AlbumImageUrl: a.AlbumImageUrl,
		}
	}
	return out
}

func toProtoAlbums(in []models.Album) []*v1.Album {
	out := make([]*v1.Album, len(in))
	for i, a := range in {
		out[i] = &v1.Album{
			Id:            a.SpotifyID,
			SpotifyId:     a.SpotifyID,
			Name:          a.Name,
			AlbumType:     a.AlbumType,
			ExternalUrls:  a.ExternalUrls,
			AlbumImageUrl: a.AlbumImageUrl,
			Genres:        a.Genres,
			Artists:       toProtoArtists(a.Artists),
		}
	}
	return out
}
