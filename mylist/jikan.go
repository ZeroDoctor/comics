package mylist

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type jResult struct {
	MalID     int64     `json:"mal_id"`
	URL       string    `json:"url"`
	ImageURL  string    `json:"image_url"`
	Title     string    `json:"title"`
	Synopsis  string    `json:"synopsis"`
	Type      string    `json:"type"`
	Chapters  int       `json:"chapters"`
	Volumes   int       `json:"volumes"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type QMediaJikanResp struct {
	Results []jResult `json:"results"`
}

func SearchMangaJikan(query string) (QMediaJikanResp, error) {
	url := "https://api.jikan.moe/v3/search/manga?q=" + query + "&page=1"
	var media QMediaJikanResp

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return media, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return media, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return media, err
	}

	err = json.Unmarshal(body, &media)
	if err != nil {
		return media, err
	}

	return media, nil
}
