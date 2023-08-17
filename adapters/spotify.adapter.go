package adapters

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type SpotifyUserToken struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Refresh_token string `json:"refresh_token"`
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

func (a *Adapters) GetSpotifyUserAccessToken(code string, client_id string, client_secret string) error {
	fmt.Println("GetSpotifyUserAccessToken")
	// Retrieve token from api
	urlReguest := "https://accounts.spotify.com/api/token"

	parm := url.Values{}
	parm.Add("code", code)
	parm.Add("redirect_uri", "http://localhost:4000/v1/spotify/callback")
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
	}
	// req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(client_id, client_secret)

	fmt.Println(req)
	fmt.Println("Req body:", req.Body)

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
	fmt.Println("GetSpotifyUserAccessToken response:", string(body))

	var userToken SpotifyUserToken
	if err := json.Unmarshal(body, &userToken); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return err
	}

	fmt.Println("JSON response", userToken)

	GetSpotifyUserData(userToken.Access_token)
	GetSpotifyUserFollowedArtists(userToken.Access_token)

	return nil

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

func GetSpotifyUserFollowedArtists(userToken string) error {
	fmt.Println("**************************************************************************")
	fmt.Println("GetSpotifyUserFollowedArtists: ", userToken)

	// Retrieve token from api
	urlReguest := "https://api.spotify.com/v1/me/following?type=artist"

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
