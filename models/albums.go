package models

type Album struct {
	ID     			string	`json:"id"`
	SpotifyID		string	`json:"spotify_id"`
	Name			string	`json:"name"`
	AlbumType 		string	`json:"album_type"`
	ExternalUrls	string 	`json:"external_urls"`
	AlbumImageUrl	string 	`json:"album_image_urls"`
	Genres			string	`json:"genres"`
	// Artists			[]Artist	`json:"artists"`
}

type AlbumList []Album

func (a AlbumList) Len() int {
    return len(a)
}

func (a AlbumList) Less(i, j int) bool {
    return a[j].Name > a[i].Name
}

func (a AlbumList) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}