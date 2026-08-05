package adapters

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	// "log"
	"recordstore-go/models"
)

// TODO: The structs here were created with AI and require refactoring to be more readable,
// Speed is great though with the retrieval.

type YoutubeUserToken struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Refresh_token string `json:"refresh_token"`
}

type YoutubeArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type YoutubeThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type YoutubeThumbnails struct {
	Thumbnails []YoutubeThumbnail `json:"thumbnails"`
}

func firstThumbnailURL(thumbnails []YoutubeThumbnail) string {
	if len(thumbnails) == 0 {
		return ""
	}
	return thumbnails[len(thumbnails)-1].URL
}

type YoutubeAlbum struct {
	ID           string            `json:"id"`
	Name         string            `json:"title"`
	AlbumType    string            `json:"type"`
	ExternalUrls string            `json:"externalUrls"`
	Thumbnails   YoutubeThumbnails `json:"thumbnails"`
	Artists      []YoutubeArtist   `json:"artists"`
}

type youtubeMusicItem struct {
	MusicTwoRowItemRenderer struct {
		Title              youtubeMusicLabel     `json:"title"`
		Subtitle           youtubeMusicLabel     `json:"subtitle"`
		ThumbnailRenderer  YoutubeMusicThumbnail `json:"thumbnailRenderer"`
		NavigationEndpoint struct {
			WatchEndpoint struct {
				VideoID string `json:"videoId"`
			} `json:"watchEndpoint"`
			BrowseEndpoint struct {
				BrowseID string `json:"browseId"`
			} `json:"browseEndpoint"`
		} `json:"navigationEndpoint"`
	} `json:"musicTwoRowItemRenderer"`
}

// youtubeMusicLabel decodes either {"text": "..."} or {"runs": [{"text": "..."}, ...]}.
// We only care about the first run's text; formatting is discarded.
type youtubeMusicLabel struct {
	Text       string            `json:"text,omitempty"`
	SimpleText string            `json:"simpleText,omitempty"`
	Runs       []youtubeMusicRun `json:"runs,omitempty"`
}

func (l youtubeMusicLabel) first() string {
	if l.Text != "" {
		return l.Text
	}
	if l.SimpleText != "" {
		return l.SimpleText
	}
	if len(l.Runs) > 0 {
		return l.Runs[0].Text
	}
	return ""
}

type youtubeMusicRun struct {
	Text string `json:"text"`
}

type YoutubeMusicThumbnail struct {
	MusicThumbnailRenderer struct {
		Thumbnail YoutubeThumbnails `json:"thumbnail"`
	} `json:"musicThumbnailRenderer"`
}

type youtubeMusicShelf struct {
	MusicShelfRenderer struct {
		Title    youtubeMusicLabel  `json:"title"`
		Contents []youtubeMusicItem `json:"contents"`
	} `json:"musicShelfRenderer"`
	MusicCarouselShelfRenderer struct {
		Header   youtubeMusicLabel  `json:"header"`
		Contents []youtubeMusicItem `json:"contents"`
	} `json:"musicCarouselShelfRenderer"`
}

func (s *youtubeMusicShelf) Items() []youtubeMusicItem {
	if len(s.MusicShelfRenderer.Contents) > 0 {
		return s.MusicShelfRenderer.Contents
	}
	return s.MusicCarouselShelfRenderer.Contents
}

type youtubeMusicTile struct {
	TileRenderer struct {
		Header struct {
			TileHeaderRenderer struct {
				Thumbnail YoutubeThumbnails `json:"thumbnail"`
			} `json:"tileHeaderRenderer"`
		} `json:"header"`
		Metadata struct {
			TileMetadataRenderer struct {
				Title struct {
					Runs []youtubeMusicRun `json:"runs"`
				} `json:"title"`
				Lines []struct {
					LineRenderer struct {
						Items []struct {
							LineItemRenderer struct {
								Text youtubeMusicLabel `json:"text"`
							} `json:"lineItemRenderer"`
						} `json:"items"`
					} `json:"lineRenderer"`
				} `json:"lines"`
			} `json:"tileMetadataRenderer"`
		} `json:"metadata"`
		OnSelectCommand struct {
			BrowseEndpoint struct {
				BrowseID string `json:"browseId"`
			} `json:"browseEndpoint"`
		} `json:"onSelectCommand"`
	} `json:"tileRenderer"`
}

type youtubeMusicLibraryResponse struct {
	Contents struct {
		SingleColumnBrowseResultsContent struct {
			Contents []youtubeMusicShelf `json:"contents"`
		} `json:"singleColumnBrowseResultsContent"`
		TabbedResultsContent struct {
			Contents []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []youtubeMusicShelf `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"contents"`
		} `json:"tabbedResultsContent"`
		TvBrowseRenderer struct {
			Content struct {
				TvSecondaryNavRenderer struct {
					Sections []struct {
						TvSecondaryNavSectionRenderer struct {
							Tabs []struct {
								TabRenderer struct {
									Content struct {
										TvSurfaceContentRenderer struct {
											Content struct {
												GridRenderer struct {
													Items []youtubeMusicTile `json:"items"`
												} `json:"gridRenderer"`
											} `json:"content"`
										} `json:"tvSurfaceContentRenderer"`
									} `json:"content"`
								} `json:"tabRenderer"`
							} `json:"tabs"`
						} `json:"tvSecondaryNavSectionRenderer"`
					} `json:"sections"`
				} `json:"tvSecondaryNavRenderer"`
			} `json:"content"`
		} `json:"tvBrowseRenderer"`
	} `json:"contents"`
}

// GetYoutubeUserAccessToken exchanges a Google OAuth2 authorization code for an
// access token and refresh token by POSTing to https://oauth2.googleapis.com/token.
// It returns the parsed YoutubeUserToken (containing access_token, refresh_token,
// token_type, etc.) and any error from the HTTP exchange or JSON decoding.
func (a *Adapters) GetYoutubeUserAccessToken(code string, client_id string, client_secret string, redirectHost string) (error, YoutubeUserToken) {
	fmt.Println("GetYoutubeUserAccessToken", redirectHost)

	// Retrieve token from api
	urlReguest := "https://oauth2.googleapis.com/token"

	parm := url.Values{}
	parm.Add("code", code)
	parm.Add("client_id", client_id)
	parm.Add("client_secret", client_secret)

	parm.Add("redirect_uri", redirectHost)
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
		return err, YoutubeUserToken{}
	}
	// req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(client_id, client_secret)

	fmt.Println("Req body:", req.Body)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return err, YoutubeUserToken{}
	}
	fmt.Println(resp.Header)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return err, YoutubeUserToken{}
	}
	// fmt.Println("GetYoutubeUserAccessToken response:", string(body))

	var userToken YoutubeUserToken
	if err := json.Unmarshal(body, &userToken); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return err, YoutubeUserToken{}
	}

	// fmt.Println("GetYoutubeUserAccessToken JSON response", userToken.Access_token)

	return nil, userToken
}

func (a *Adapters) GetYoutubeUserSavedAlbums(userToken string) (error, []models.Album) {
	fmt.Println("**************************************************************************")
	fmt.Println("GetYoutubeUserSavedAlbums Enter")
	fmt.Println("**************************************************************************")

	var albums []models.Album

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	clientVersion := "7.20260724.01.00"
	requestBody, err := json.Marshal(map[string]interface{}{
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"clientName":    "TVHTML5",
				"clientVersion": clientVersion,
				"hl":            "en",
				"gl":            "US",
			},
		},
		"browseId": "FEmusic_liked_albums",
	})
	if err != nil {
		fmt.Println("GetYoutubeUserSavedAlbums Error: ", err)
		return err, albums
	}

	u, _ := url.Parse("https://music.youtube.com")
	cookieJar.SetCookies(u, []*http.Cookie{{Name: "SOCS", Value: "CAI"}})

	visitorID, visitorErr := a.fetchYouTubeVisitorID(client)
	if visitorErr != nil {
		fmt.Println("GetYoutubeUserSavedAlbums visitor fetch failed:", visitorErr)
	}

	req, err := http.NewRequest("POST", "https://music.youtube.com/youtubei/v1/browse?alt=json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println("GetYoutubeUserSavedAlbums Error: ", err)
		return err, albums
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Authorization", "Bearer "+userToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", "https://music.youtube.com")
	req.Header.Add("Referer", "https://music.youtube.com/")
	if visitorID != "" {
		req.Header.Add("X-Goog-Visitor-Id", visitorID)
	}

	// fmt.Println("GetYoutubeUserSavedAlbums request header:", req.Header)
	// fmt.Println("GetYoutubeUserSavedAlbums request body:", string(requestBody))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("GetYoutubeUserSavedAlbums Error: ", err)
		return err, albums
	}
	defer resp.Body.Close()

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("GetYoutubeUserSavedAlbums Error creating gzip reader:", err)
			return err, albums
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("GetYoutubeUserSavedAlbums Error: ", err)
		return err, albums
	}
	// fmt.Println("GetYoutubeUserSavedAlbums body length:", len(body))
	// fmt.Println("GetYoutubeUserSavedAlbums header:", resp.Header)
	// fmt.Println("GetYoutubeUserSavedAlbums body:", string(body))

	var parsed youtubeMusicLibraryResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		fmt.Println("GetYoutubeUserSavedAlbums Error: Can not unmarshal JSON")
		return err, albums
	}

	for _, tab := range parsed.Contents.TabbedResultsContent.Contents {
		for _, shelf := range tab.TabRenderer.Content.SectionListRenderer.Contents {
			for _, item := range shelf.Items() {
				r := item.MusicTwoRowItemRenderer
				if r.NavigationEndpoint.BrowseEndpoint.BrowseID == "" {
					continue
				}

				albumImage := ""
				if len(r.ThumbnailRenderer.MusicThumbnailRenderer.Thumbnail.Thumbnails) > 0 {
					albumImage = r.ThumbnailRenderer.MusicThumbnailRenderer.Thumbnail.Thumbnails[0].URL
				}

				albumName := r.Title.first()
				var ytArtists []YoutubeArtist
				if subtitle := r.Subtitle.first(); subtitle != "" {
					ytArtists = append(ytArtists, YoutubeArtist{Name: subtitle})
				}

				artists := make([]models.Artist, 0, len(ytArtists))
				for _, ya := range ytArtists {
					artists = append(artists, models.Artist{
						Name: ya.Name,
					})
				}

				curAlbum := models.Album{
					Name:          albumName,
					AlbumType:     "album",
					ExternalUrls:  "https://music.youtube.com/browse/" + r.NavigationEndpoint.BrowseEndpoint.BrowseID,
					AlbumImageUrl: albumImage,
					Artists:       artists,
				}
				albums = append(albums, curAlbum)
			}
		}
	}

	for _, section := range parsed.Contents.TvBrowseRenderer.Content.TvSecondaryNavRenderer.Sections {
		for _, tab := range section.TvSecondaryNavSectionRenderer.Tabs {
			for _, item := range tab.TabRenderer.Content.TvSurfaceContentRenderer.Content.GridRenderer.Items {
				r := item.TileRenderer
				if r.OnSelectCommand.BrowseEndpoint.BrowseID == "" {
					continue
				}

				albumImage := firstThumbnailURL(r.Header.TileHeaderRenderer.Thumbnail.Thumbnails)

				albumName := ""
				if len(r.Metadata.TileMetadataRenderer.Title.Runs) > 0 {
					albumName = r.Metadata.TileMetadataRenderer.Title.Runs[0].Text
				}

				var ytArtists []YoutubeArtist
				if len(r.Metadata.TileMetadataRenderer.Lines) > 0 && len(r.Metadata.TileMetadataRenderer.Lines[0].LineRenderer.Items) > 0 {
					ytArtists = append(ytArtists, YoutubeArtist{Name: r.Metadata.TileMetadataRenderer.Lines[0].LineRenderer.Items[0].LineItemRenderer.Text.first()})
				}

				artists := make([]models.Artist, 0, len(ytArtists))
				for _, ya := range ytArtists {
					artists = append(artists, models.Artist{Name: ya.Name})
				}

				curAlbum := models.Album{
					Name:          albumName,
					AlbumType:     "album",
					ExternalUrls:  "https://music.youtube.com/browse/" + r.OnSelectCommand.BrowseEndpoint.BrowseID,
					AlbumImageUrl: albumImage,
					Artists:       artists,
				}
				albums = append(albums, curAlbum)
			}
		}
	}

	return nil, albums
}

func (a *Adapters) fetchYouTubeVisitorID(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://music.youtube.com", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Referer", "https://music.youtube.com/")
	req.Header.Add("Origin", "https://music.youtube.com")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"VISITOR_DATA"\s*:\s*"([^"]+)"`)
	match := re.FindSubmatch(body)
	if len(match) < 2 {
		return "", fmt.Errorf("missing VISITOR_DATA in YouTube page")
	}

	return string(match[1]), nil
}

func (a *Adapters) GetGoogleUserDisplayName(userToken string) (string, error) {
	fmt.Println("GetGoogleUserDisplayName Enter", userToken)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", "https://openidconnect.googleapis.com/v1/userinfo", strings.NewReader(""))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+userToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var profile struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := json.Unmarshal(body, &profile); err != nil {
		return "", err
	}

	fmt.Println("GetGoogleUserDisplayName profile:", profile.Name, profile.Email, profile.Verified)
	// Use email as username
	return profile.Email, nil
}
