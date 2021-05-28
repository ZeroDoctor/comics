package mylist

import (
	"github.com/zerodoctor/graphql"
)

type QMediaAnilist struct {
	Search string
	Page   int
}

type QMediaAnilistResp struct {
	Page aPage `json:"Page"`
}

func SearchMangaAnilist(query QMediaAnilist) (QMediaAnilistResp, error) {
	url := "https://graphql.anilist.co"

	req := graphql.NewRequest(`
query($page: Int!, $search: String!){ 
	Page(page: $page) {
		pageInfo {
			total
			perPage
			currentPage
			lastPage
			hasNextPage
		}
		media(search: $search, type: MANGA) {
			id
			type
			source
			countryOfOrigin
			genres
			synonyms
			startDate {
				year
				month
				day
			}
			endDate {
				year
				month
				day
			}
			title {
				romaji
				english
				userPreferred
			}
			isAdult
		}
	}
}`)

	req.Var("page", query.Page)
	req.Var("search", query.Search)

	var media QMediaAnilistResp
	err := newApi().FetchAnilist(url, req, &media)

	return media, err
}
