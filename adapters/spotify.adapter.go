package adapters

import (
	"crypto/tls"
   	"fmt"
   	"io/ioutil"
   	"net/http"
   	"net/http/cookiejar"
   	"net/url"
   	"strings"
)

// OpenSpotifyConnection gets a bearer to use for this session
func (a *Adapters) OpenSpotifyConnection() error {
		// Retrieve token from api
		urlReguest := a.baseUrl

		parm := url.Values{}
		parm.Add("grant_type", "client_credentials")
		parm.Add("client_id", "5c37b6f4c90143908b11d9e1727db5e7")
		parm.Add("client_secret", "0a5f8bc725a94eb7aae28803adc108d1")

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

		// resp, err := http.Get(urlReguest)
		// if err != nil {1
		// 	log.Fatal("worker error ", urlReguest)
		// 	return
		// }

}

func (a *Adapters) OpenSpotifyUserConnection() error {
		// Retrieve token from api
		urlReguest := "https://accounts.spotify.com/authorize"

		parm := url.Values{}
		parm.Add("response_type", "code")
		parm.Add("client_id", "5c37b6f4c90143908b11d9e1727db5e7")
		parm.Add("scope", "scope")
		parm.Add("redirect_uri", "http:/localhost:8181")
		parm.Add("state", "state")

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

		// resp, err := http.Get(urlReguest)
		// if err != nil {1
		// 	log.Fatal("worker error ", urlReguest)
		// 	return
		// }

}