package main

import (
	"fmt"
	"time"

	"github.com/zerodoctor/comics/mylist"
	"github.com/zerodoctor/comics/ui"
)

func main() {

	/* op := webtoons.Option{
		TitleNum: "98",
		Genre:    "GENRE", // default GENRE
		Title:    "TITLE", // default TITLE
		Start:    1,       // default 1
		End:      -1,      // default -1
		Workers:  5,       // default 10
		Verbose:  false,   // default false
	}

	webtoons.Scrap(op)
	*/

	// mylist.SearchMangaJikan("One%20pun")

	ui.Start()

	return

	query := mylist.QMediaAnilist{
		Search: "One punch",
		Page:   1,
	}

	manga, _ := mylist.SearchMangaAnilist(query)
	fmt.Printf("%+v\n", manga)

	time.Sleep(time.Millisecond * 750)
}
