package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Album struct {
	ID            uint       `json:"id"`
	SpotifyID     string     `json:"spotify_id"`
	Name          string     `json:"name"`
	AlbumType     string     `json:"album_type"`
	ExternalUrls  string     `json:"external_urls"`
	AlbumImageUrl string     `json:"album_image_urls"`
	Genres        string     `json:"genres"`
	Artists       ArtistList `json:"artists"`
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

func (a Album) Insert(db *gorm.DB) error {
	result := db.Create(&a)
	if result.Error != nil {
		fmt.Println("InsertArtist error", result.Error)
		return result.Error
	}

	fmt.Println("Artist:InsertArtist Success", a)
	return nil
}

func (al AlbumList) InsertAll(db *gorm.DB) error {

	for _, album := range al {
		album.Insert(db)

		// Add user relationship

	}

	return nil
}

func (a Album) GetAlbum(db *gorm.DB) error {
	db.First(&a, a.ID)
	return nil
}

