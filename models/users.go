package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	CreatedAt           time.Time `json:"-"`
	ModifiedAt          time.Time `json:"-"`
	DeletedAt           time.Time `json:"-"`
	UserName            string    `json:"user_name"`
	SpotifySession      string    `json:"spotify_session"`
	YoutubeSession      string    `json:"youtube_session"`
	YoutubeRefreshToken string    `json:"-"`
	Analysis            string    `json:"analysis"`
}

type ConnectedProvider struct {
	UserID      uint      `gorm:"primaryKey" json:"user_id"`
	Provider    string    `gorm:"primaryKey" json:"provider"`
	ConnectedAt time.Time `json:"connected_at"`
}

type UserArtist struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	UserID   uint `json:"user_id"`
	ArtistID uint `json:"artist_id"`
	Spotify  bool `json:"spotify"`
}

type UserAlbum struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	UserID  uint `json:"user_id"`
	AlbumID uint `json:"album_id"`
	Spotify bool `json:"spotify"`
}

func (u *User) CreateUser(db *gorm.DB) error {

	result := db.Create(&u)
	if result.Error != nil {
		fmt.Println("CreateUser error", result.Error)
		return result.Error
	}

	fmt.Println("User:CreateUser Success", u)
	return nil
}

func (u *User) UpdateUser(db *gorm.DB) error {
	result := db.Save(&u)
	return result.Error
}

func GetUserByID(db *gorm.DB, id uint) (User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) GetConnectedProviders(db *gorm.DB) ([]ConnectedProvider, error) {
	var cps []ConnectedProvider
	if err := db.Where("user_id = ?", u.ID).Find(&cps).Error; err != nil {
		return cps, err
	}
	return cps, nil
}

func (u *User) CreateUserSpotifySession(db *gorm.DB, sessionID string) error {

	result := db.Create(&u)
	if result.Error != nil {
		fmt.Println("User:CreateUserSpotifySession error", result.Error)
		return result.Error
	}

	fmt.Println("User:CreateUserSpotifySession Success", u)
	return nil

}

// func (u *User) GetUser(db *gorm.DB, sessionID string) {
// 	db.First(&u, u.ID)
// }

func (u *User) GetUser(db *gorm.DB, sessionID string) {
	db.First(&u, u.ID)
}

func GetExistingSessionByUserID(db *gorm.DB, userID string) (User, error) {
	var user User
	result := db.Where("id = ?", userID).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("IsExistingSession error: RecordNotFound")
			return user, nil // Record not found, user with the given sessionID does not exist
		}
		fmt.Println("IsExistingSession error: ", result.Error)
		return user, result.Error
	}

	return user, nil // Record found, user with the given sessionID exists
}

func (u *User) UpdateUserLibrary(db *gorm.DB, md *MusicData) error {
	fmt.Println("Users:UpdateUserLibrary enter", u.ID)
	// Albums
	err := u.AddUserAlbums(db, md.Albums)
	if err != nil {
		fmt.Println("User:UpdateUserLibrary Albums error", err)
		return err
	}

	// Artists
	err = u.AddUserArtists(db, md.Artists)
	if err != nil {
		fmt.Println("User:UpdateUserLibrary Artists error", err)
		return err
	}

	return nil
}

func (u *User) AddUserArtist(db *gorm.DB, artist Artist) error {
	if artist.ID == 0 {
		artist.ID = artist.GetArtistIDBySpotifyID(db)
	}

	// Add the user artist relationship
	userArtist := UserArtist{
		UserID:   u.ID,
		ArtistID: artist.ID,
		Spotify:  true,
	}

	result := db.Create(&userArtist)
	if result.Error != nil {
		fmt.Println("AddUserArtist error", result.Error)
		return result.Error
	}

	return nil
}

func (u *User) AddUserArtists(db *gorm.DB, al ArtistList) error {
	fmt.Println("Enter User:AddUserArtists")

	for _, artist := range al {
		u.AddUserArtist(db, artist)

		// Add user relationship

	}

	return nil
}

func (u *User) GetUserArtists(db *gorm.DB) ([]Artist, error) {
	var al []Artist

	// First get user/artist associations
	var userArtists []UserArtist
	db.Where("user_id = ?", u.ID).Find(&userArtists)

	// Next fill artist list with artist data
	for _, ua := range userArtists {
		// Get artist and add to artist list
		fmt.Println("Artist ID", ua.ArtistID)
		artist := Artist{
			ID: ua.ArtistID,
		}
		artist.GetArtist(db)
		al = append(al, artist)
	}

	return al, nil
}

func (u *User) AddUserAlbum(db *gorm.DB, album Album) error {
	if album.ID == 0 {
		album.ID = album.GetAlbumIDBySpotifyID(db)
	}

	userAlbum := UserAlbum{
		UserID:  u.ID,
		AlbumID: album.ID,
		Spotify: true,
	}

	result := db.Create(&userAlbum)
	if result.Error != nil {
		fmt.Println("AddUserAlbum error", result.Error)
		return result.Error
	}

	// TODO JV Add artists to db
	// for _, artist := range album.Artists {
	// 	artist.Insert(db)
	// }

	return nil
}

func (u *User) AddUserAlbums(db *gorm.DB, al AlbumList) error {
	fmt.Println("Enter User:AddUserAlbums")

	for _, album := range al {
		u.AddUserAlbum(db, album)
	}

	return nil
}

func (u *User) GetUserAlbums(db *gorm.DB) ([]Album, error) {
	var al []Album

	// First get user/album associations
	var userAlbums []UserAlbum
	db.Where("user_id = ?", u.ID).Find(&userAlbums)

	// Next fill album list with album data
	for _, ua := range userAlbums {
		// Get album and add to album list
		album := Album{
			ID: ua.AlbumID,
		}
		album.GetAlbum(db)
		al = append(al, album)
	}

	return al, nil
}
