package adapters

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"recordstore-go/models"
)

var promptTestText = "Write a paragraph cynically describing the music tastes of someone that likes the following artists: Adrian Quesada, Alabama Shakes, Aloe Blacc,American Aquarium,Beastie Boys,Billy Bragg & Wilco,Black Country, New Road, Black Joe Lewis & The Honeybears,Black Pistol Fire,Blitzen Trapper,Bobby Jealousy,Bria,Caitlin Rose,Calexico,Dawes,Death,Dirty Projectors,Explosions In The Sky,"

func (a *Adapters) ArtistsResponseRequest(artists models.ArtistList) error {
	fmt.Println("ArtistsResponseRequest Enter")
	urlReguest := "https://api.deepai.org/api/text-generator"

	parm := url.Values{}
	parm.Add("text", promptTestText)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, err := http.NewRequest("POST", urlReguest, strings.NewReader(parm.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("api-key", "quickstart-QUdJIGlzIGNvbWluZy4uLi4K")

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
	fmt.Println("ArtistsResponseRequest", string(body))

	return nil
}
