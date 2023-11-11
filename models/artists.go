package models

type Artist struct {
	ID     			string	`json:"id"`
	SpotifyID		string	`json:"spotify_id"`
	Name			string	`json:"name"`
	ExternalUrls	string 	`json:"external_urls"`
	AlbumImageUrl	string 	`json:"album_image_urls"`
}

type ArtistList []Artist

type ArtistsAnalysis struct {
	Artists 	ArtistList	`json:"artists"` 
	Description	string		`json:"description"`
}

func (a ArtistList) Len() int {
    return len(a)
}

func (a ArtistList) Less(i, j int) bool {
    return a[j].Name > a[i].Name
}

func (a ArtistList) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}