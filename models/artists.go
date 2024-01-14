package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Artist struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	SpotifyID     string `gorm:"uniqueIndex" json:"spotify_id"`
	Name          string `json:"name"`
	ExternalUrls  string `json:"external_urls"`
	AlbumImageUrl string `json:"album_image_urls"`
}

type ArtistList []Artist

type ArtistsAnalysis struct {
	Artists     ArtistList `json:"artists"`
	Description string     `json:"description"`
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

func (a Artist) Insert(db *gorm.DB) error {
	result := db.Create(&a)
	if result.Error != nil {
		fmt.Println("InsertArtist error", result.Error)
		return result.Error
	}

	fmt.Println("Artist:InsertArtist Success", a)
	return nil
}

func (al ArtistList) InsertAll(db *gorm.DB) error {
	fmt.Println("Enter Artist:insertAll")

	for _, artist := range al {
		artist.Insert(db)
	}

	return nil
}

func (a Artist) GetArtist(db *gorm.DB) error {
	db.First(&a, a.ID)
	return nil
}
