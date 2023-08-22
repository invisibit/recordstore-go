package models

type Artist struct {
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Genres []string `json:"genres"`
	Href   string   `json:"href"`
	ID     string   `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	Type       string `json:"type"`
	URI        string `json:"uri"`
}


type Artists struct {
	// Href	string 		`json:"href"`
	// Limit	int 		`json:"limit"`
	// Next	string 		`json:"next"`	
	// Cursors	Cursors		`json:"cursors"`
	// Total	string 		`json:"total"`
	Items	[]Artist 	`json:"items"`
}

type Followed struct {
	Artists	Artists		`json:"artists"`
}

