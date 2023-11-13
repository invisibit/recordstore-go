package models

type MusicData struct {
	Albums 		AlbumList	`json:"albums"` 
	Artists 	ArtistList	`json:"artists"` 
	Analysis	string		`json:"analysis"`
}

