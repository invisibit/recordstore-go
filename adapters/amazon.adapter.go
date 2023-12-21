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
	// "sort"
	"strings"

	// "recordstore-go/models"
)

type AmazonUserToken struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Refresh_token string `json:"refresh_token"`
}

func (a *Adapters) GetAmazonUserAccessToken(code string, client_id string, client_secret string, redirectHost string) (error, string) {
	fmt.Println("GetAmazonUserAccessToken")
	fmt.Println("GetAmazonUserAccessToken", redirectHost)

	// Retrieve token from api
	urlReguest := "https://api.amazon.com/auth/o2/token"

	parm := url.Values{}
	parm.Add("code", code)

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
		return err, ""
	}
	// req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(client_id, client_secret)

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

	var userToken AmazonUserToken
	if err := json.Unmarshal(body, &userToken); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return err, ""
	}

	fmt.Println("JSON response", userToken.Access_token)

	return nil, userToken.Access_token
}

