package main

import (
	"context"
	"fmt"
	"io"

	"connectrpc.com/connect"

	"recordstore-go/adapters"
	v1 "recordstore-go/gen/recordstore/v1"
	"recordstore-go/models"
)

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

func (app *application) GetSavedAlbums(ctx context.Context, req *connect.Request[v1.GetSavedAlbumsRequest]) (*connect.Response[v1.GetSavedAlbumsResponse], error) {
	adapter := adapters.NewAdapter("")
	err, albums := adapter.GetSpotifyUserSavedAlbums(req.Msg.SptfySession)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.GetSavedAlbumsResponse{
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
