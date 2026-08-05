package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Album struct {
	ID                   uint       `json:"id"`
	SpotifyID            string     `json:"spotify_id"`
	Name                 string     `json:"name"`
	AlbumType            string     `json:"album_type"`
	TotalTracks          int        `json:"total_tracks"`
	AvailableMarkets     string     `json:"available_markets"`
	ExternalUrls         string     `json:"external_urls"`
	Href                 string     `json:"href"`
	AlbumImageUrl        string     `json:"album_image_urls"`
	ReleaseDate          string     `json:"release_date"`
	ReleaseDatePrecision string     `json:"release_date_precision"`
	Restrictions         string     `json:"restrictions"`
	SpotifyURI           string     `json:"spotify_uri"`
	Genres               string     `json:"genres"` // Deprecated
	Artists              ArtistList `gorm:"-" json:"artists"`
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
	// Check if exists
	var artist Artist
	db.Where("spotify_id = ?", a.SpotifyID).Find(&artist)
	if artist.ID != 0 {
		a.ID = artist.ID
		return nil
	}

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

func (a *Album) GetAlbum(db *gorm.DB) error {
	db.First(&a, a.ID)
	return nil
}

func (a *Album) GetAlbumIDBySpotifyID(db *gorm.DB) uint {
	var album Album
	db.Where("spotify_id = ?", a.SpotifyID).Find(&album)

	return album.ID
}
