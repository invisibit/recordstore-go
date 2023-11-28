package adapters

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strings"

	"recordstore-go/models"
)

type SpotifyUserToken struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Refresh_token string `json:"refresh_token"`
}

type SpotifyArtist struct {
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	// Followers struct {
	// 	Href  string `json:"href"`
	// 	Total int    `json:"total"`
	// } `json:"followers"`
	// Genres []string `json:"genres"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name string `json:"name"`
	// Popularity int    `json:"popularity"`
	// Type       string `json:"type"`
	// URI        string `json:"uri"`
}

type SpotifyArtists struct {
	// Href	string 		`json:"href"`
	// Limit	int 		`json:"limit"`
	Items   []SpotifyArtist `json:"items"`
	Next    string          `json:"next"`
	Cursors struct {
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"cursors"`
	Total int `json:"total"`
}

type SpotifyFollowed struct {
	Artists SpotifyArtists `json:"artists"`
}

type SpotifyAlbum struct {
	AlbumType    string `json:"album_type"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type SpotifyItem struct {
	AddedAt string       `json:"added_at"`
	Album   SpotifyAlbum `json:"album"`
}

type SpotifyAlbums struct {
	// Href	string 		`json:"href"`
	// Limit	int 		`json:"limit"`
	Items []SpotifyItem `json:"items"`
	Next  string        `json:"next"`
	Total int           `json:"total"`
}

type SpotifySavedAlbums struct {
	Albums SpotifyAlbums `json:"albums"`
}

// OpenSpotifyConnection gets a bearer to use for this session
func (a *Adapters) OpenSpotifyConnection(client_id string, client_secret string) error {
	// Retrieve token from api
	urlReguest := a.baseUrl

	parm := url.Values{}
	parm.Add("grant_type", "client_credentials")
	parm.Add("client_id", client_id)
	parm.Add("client_secret", client_secret)

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
	}
	// fmt.Println(resp.Header)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	fmt.Println(string(body))

	return nil

}

func (a *Adapters) GetSpotifyUserAccessToken(code string, client_id string, client_secret string) (error, string) {
	fmt.Println("GetSpotifyUserAccessToken")
	// Retrieve token from api
	urlReguest := "https://accounts.spotify.com/api/token"

	parm := url.Values{}
	parm.Add("code", code)
	parm.Add("redirect_uri", "https://recordstore-go-344gqgcrvq-uc.a.run.app/v1/spotify/callback")
	parm.Add("grant_type", "authorization_code")

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
	if err != nil {
		fmt.Println("New Request Error: ", err)
		return err, ""
	}
	// req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(client_id, client_secret)

	fmt.Println(req)
	fmt.Println("Req body:", req.Body)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return err, ""
	}
	fmt.Println(resp.Header)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return err, ""
	}
	fmt.Println("GetSpotifyUserAccessToken response:", string(body))

	var userToken SpotifyUserToken
	if err := json.Unmarshal(body, &userToken); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return err, ""
	}

	fmt.Println("JSON response", userToken.Access_token)

	return nil, userToken.Access_token

	// GetSpotifyUserData(userToken.Access_token)
	// err, followedArtists = GetSpotifyUserFollowedArtists(userToken.Access_token)

	// return nil, followedArtists

	// resp, err := http.Get(urlReguest)
	// if err != nil {1
	// 	log.Fatal("worker error ", urlReguest)
	// 	return
	// }

}

func GetSpotifyUserData(userToken string) error {
	fmt.Println("GetSpotifyUserData: ", userToken)

	// Retrieve token from api
	urlReguest := "https://api.spotify.com/v1/me"

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, err := http.NewRequest("GET", urlReguest, strings.NewReader(""))
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+userToken)

	fmt.Println(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(resp.Header)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	fmt.Println(string(body))

	return nil
}

func (a *Adapters) GetSpotifyUserFollowedArtists(userToken string) (error, []models.Artist) {
	fmt.Println("**************************************************************************")
	fmt.Println("GetSpotifyUserFollowedArtists: ")

	var followedArtists SpotifyFollowed
	var artists models.ArtistList
	// Retrieve token from api
	urlRequest := "https://api.spotify.com/v1/me/following?type=artist"

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	var isFinished = false
	for !isFinished {
		req, err := http.NewRequest("GET", urlRequest, strings.NewReader(""))
		req.Header.Add("Authorization", "Bearer "+userToken)
		// fmt.Println("GetSpotifyUserFollowedArtists", req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Println(resp.Header)
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

		fmt.Println(resp.Body)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return err, artists
		}

		// Parse []byte to the go struct pointer
		if err := json.Unmarshal(body, &followedArtists); err != nil {
			fmt.Println("Can not unmarshal JSON")
			return err, artists
		}

		// if len(artists) > 39 {
		// 	fmt.Println("followedArtists.Artists", followedArtists.Artists)
		// }
		// Convert return Artists to Artist model
		for _, artist := range followedArtists.Artists.Items {
			// Check if there is an associated image for the artist
			artistImage := ""
			if len(artist.Images) > 0 {
				artistImage = artist.Images[0].URL

			}

			curArtist := models.Artist{
				ID:            artist.ID,
				Name:          artist.Name,
				ExternalUrls:  artist.ExternalUrls.Spotify,
				AlbumImageUrl: artistImage, // decide which one
			}
			artists = append(artists, curArtist)
		}

		urlRequest = followedArtists.Artists.Next

		if len(artists) >= followedArtists.Artists.Total {
			isFinished = true
		}
	}

	// Sort by artist name
	sort.Sort(models.ArtistList(artists))

	return nil, artists
}

func (a *Adapters) GetSpotifyUserSavedAlbums(userToken string) (error, []models.Album) {
	fmt.Println("**************************************************************************")
	fmt.Println("GetSpotifyUserSavedAlbums Enter")

	var savedAlbums SpotifyAlbums
	var albums []models.Album

	// Retrieve token from api
	urlRequest := "https://api.spotify.com/v1/me/albums"

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	var isFinished = false
	for !isFinished {
		req, err := http.NewRequest("GET", urlRequest, strings.NewReader(""))
		req.Header.Add("Authorization", "Bearer "+userToken)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("%s", err)
			return err, albums
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return err, albums
		}

		// Parse []byte to the go struct pointer
		if err := json.Unmarshal(body, &savedAlbums); err != nil {
			fmt.Println("Can not unmarshal JSON")
			return err, albums
		}

		// if len(albums) > 39 {
		// 	fmt.Println("savedAlbums.Artists", savedAlbums.Albums)
		// }
		// fmt.Println("Saved:", savedAlbums)
		// Convert return Artists to Artist model
		// fmt.Println("GetSpotifyUserSavedAlbums number albums", len(savedAlbums.Items))
		for _, item := range savedAlbums.Items {
			// Check if there is an associated image for the artist
			albumImage := ""
			if len(item.Album.Images) > 0 {
				albumImage = item.Album.Images[0].URL

			}

			var genres = ""
			for _, genre := range item.Album.Genres {
				genres += genre + " "
			}

			curAlbum := models.Album{
				ID:            item.Album.ID,
				Name:          item.Album.Name,
				AlbumType:     item.Album.AlbumType,
				ExternalUrls:  item.Album.ExternalUrls.Spotify,
				AlbumImageUrl: albumImage, // decide which one
				Genres:        genres,
			}
			albums = append(albums, curAlbum)
		}

		urlRequest = savedAlbums.Next

		if len(albums) >= savedAlbums.Total {
			isFinished = true
		}
	}

	return nil, albums
}
