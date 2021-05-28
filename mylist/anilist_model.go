package mylist

import (
	"encoding/json"
	"time"
)

type AnilistTime struct{ time.Time }

func (a *AnilistTime) UnmarshalJSON(data []byte) error {
	var Date struct {
		Year  *int `json:"year"`
		Month *int `json:"month"`
		Day   *int `json:"day"`
	}

	if err := json.Unmarshal(data, &Date); err != nil {
		return err
	}

	year := 0
	month := 1
	day := 1

	if Date.Year != nil {
		year = *Date.Year
	}
	if Date.Month != nil {
		month = *Date.Month
	}
	if Date.Day != nil {
		day = *Date.Day
	}

	a.Time = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return nil
}

type aPageInfo struct {
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	Total       int  `json:"total"`
	HasNextPage bool `json:"hasNextPage"`
}

type aTitle struct {
	Romaji        string `json:"romaji"`
	English       string `json:"english"`
	UserPreferred string `json:"userPreferred"`
}

type aMedia struct {
	ID              int64       `json:"id"`
	Type            string      `json:"type"`
	Source          string      `json:"source"`
	Genres          []string    `json:"genres"`
	CountryOfOrigin string      `json:"countryOfOrigin"`
	Synonyms        []string    `json:"synonyms"`
	StartDate       AnilistTime `json:"startDate"`
	EndDate         AnilistTime `json:"endDate"`
	Title           aTitle      `json:"title"`
	IsAdult         bool        `json:"isAdult"`
}

type aPage struct {
	PageInfo aPageInfo `json:"pageInfo"`
	Media    []aMedia  `json:"media"`
}
