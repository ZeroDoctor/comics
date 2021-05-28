package webtoons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/zerodoctor/comics/util"
)

// scrap main function for scrapping webtoons
func Scrap(op Option) {
	// * note: if webtoons update there url schema, we would have to figure out this all over again
	list := "https://www.webtoons.com/en/" + op.Genre + "/" + op.Title + "/list?title_no=" + op.TitleNum
	url := "http://www.webtoons.com/en/" + op.Genre + "/" + op.Title + "/CHAPTER/viewer?title_no=" + op.TitleNum + "&episode_no="

	info := make(chan Info, 1)

	// Log url here?
	scrapInfo(list, op, info)

	comic := <-info
	url += comic.Episode
	createFolder(comic)
	// Log result here?

	wait := make(chan bool, 1)
	episodeMap := make(map[string]string, op.End-op.Start)
	scrapComic(url, comic, op, episodeMap, wait)
	<-wait

	var procs []util.Process
	for url, name := range episodeMap {
		p := util.Process{
			Name: name,
			Args: []interface{}{url, &comic, &op, episodeMap},
			Fn: func(p util.Process) {
				scrapEpisode(p.Args[0].(string), p.Args[1].(*Info), p.Args[2].(*Option), p.Args[3].(map[string]string))
			},
		}

		procs = append(procs, p)
	}

	util.WorkerPool(procs, int(op.Workers))
}

func createFolder(comic Info) {
	var err error
	if _, err = os.Stat("./" + comic.Title); os.IsNotExist(err) {
		err := os.Mkdir("./"+comic.Title, 0755)
		if err != nil {
			log.Println("failed to create folder:", err.Error())
			os.Exit(1)
		}
	}

	if err != nil && !os.IsNotExist(err) {
		log.Println("failed to create folder:", err.Error())
		os.Exit(1)
	}

	if _, err = os.Stat("./" + comic.Title + "/about.json"); os.IsNotExist(err) {
		data, err := json.MarshalIndent(comic, "", "  ")
		if err != nil {
			log.Println("failed to marshal comic infos:", err.Error())
			os.Exit(1)
		}

		err = ioutil.WriteFile("./"+comic.Title+"/about.json", data, 0755)
		if err != nil {
			log.Println("failed to write to file:", err.Error())
			os.Exit(1)
		}
	}

	if err != nil && !os.IsNotExist(err) {
		log.Println("failed to create about file:", err.Error())
		os.Exit(1)
	}
}

func scrapInfo(list string, op Option, info chan Info) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{list},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			var comic Info
			var err error
			var ok bool

			comic.Title = util.CleanString(r.HTMLDoc.Find(".info").Find(".subj").Text())
			comic.Subscriber = r.HTMLDoc.Find(".grade_area").Find("span.ico_subscribe + em").Text()
			comic.Rating = r.HTMLDoc.Find("#_starScoreAverage").Text()
			comic.Summary = r.HTMLDoc.Find("#_asideDetail > p.summary").Text()
			comic.Episode, ok = r.HTMLDoc.Find("#_listUl > li").Attr("data-episode-no")
			if !ok {
				log.Println("failed to parse latest episode url:", err.Error())
				os.Exit(1)
			}
			endStr := r.HTMLDoc.Find("#_listUl > li").Find(".tx").Text()
			endStr = strings.Split(endStr, "#")[1]

			end := strings.Replace(endStr, "#", "", -1)
			comic.End, err = strconv.Atoi(end)
			if err != nil {
				errStr := fmt.Sprintln("failed to parse latest episode number:", err.Error())
				util.AddLog(errStr, comic.Title)
			}

			var prefixes []string
			var creators []string
			r.HTMLDoc.Find("div._authorInfoLayer div._authorInnerContent").Find("p.by").Each(
				func(_ int, s *goquery.Selection) {
					prefixes = append(prefixes, s.Text())
				},
			)

			r.HTMLDoc.Find("div._authorInfoLayer div._authorInnerContent").Find("h3.title").Each(
				func(_ int, s *goquery.Selection) {
					creators = append(creators, s.Text())
				},
			)

			for i := range prefixes {
				comic.Creators = append(comic.Creators, prefixes[i]+": "+creators[i])
			}
			info <- comic
		},
		LogDisabled: !op.Verbose,
	}).Start()
}

// scrapComic for episode list
func scrapComic(url string, comic Info, op Option, episodeMap map[string]string, wait chan bool) {
	if op.End <= 1 {
		op.End = comic.End
	}
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("#topEpisodeList").Find("div.episode_cont").Find("li").EachWithBreak(
				func(i int, s *goquery.Selection) bool {
					num := i + 1
					if num < op.Start {
						return true
					} else if op.End != -1 && num > op.End {
						return false
					}
					next, _ := s.Find("a").Attr("href")
					title, _ := s.Find("img").Attr("alt")
					episodeMap[next] = fmt.Sprintf("[%d]", num) + title
					return true
				},
			)

			wait <- false
		},
		LogDisabled: !op.Verbose,
	}).Start()
}

func scrapEpisode(urlStr string, comic *Info, op *Option, episodeMap map[string]string) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{urlStr},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("recovered from painc:", r)
					os.Exit(1)
				}
			}()

			width := 8.33
			height := 13.33

			var panels []util.Panel
			r.HTMLDoc.Find("#_imageList").Find("img").Each(
				func(counter int, s *goquery.Selection) {
					// find panel image url
					href, ok := s.Attr("data-url")
					if !ok {
						errStr := fmt.Sprintf("failed read data-url (%d)\n", counter)
						util.AddLog(errStr, comic.Title)
						return
					}

					url, err := r.JoinURL(href)
					if err != nil {
						errStr := fmt.Sprintf("failed parse image url (%d): %s\n", counter, err.Error())
						util.AddLog(errStr, comic.Title)
						return
					}

					// create get request with important header
					req := &http.Request{
						Method: "GET",
						Header: http.Header(map[string][]string{
							// * note: super important header. if changed, thing will become a lot harder
							"Referer": {"http://www.webtoons.com"},
						}),
						URL: url,
					}

					// send request
					resp, err := g.Client.Do(req)
					if err != nil {
						errStr := fmt.Sprintln("failed request:", err.Error())
						util.AddLog(errStr, comic.Title)
					}

					// handle response
					data, err := util.ReadImageFromResp(resp)
					if err != nil {
						errStr := fmt.Sprintln("failed read resp image:", err.Error())
						util.AddLog(errStr, comic.Title)
						return
					}

					imageType := resp.Header["Content-Type"][0][len("image/"):]

					imgWidth, ok := s.Attr("width")
					if ok {
						w, err := strconv.ParseFloat(imgWidth, 64)
						if err != nil {
							errStr := fmt.Sprintln("failed to parse width:", err.Error())
							util.AddLog(errStr, comic.Title)
						} else {
							width = float64(w+15) * 0.0104166667
						}
					}

					imgHeight, ok := s.Attr("height")
					if ok {
						h, err := strconv.ParseFloat(imgHeight, 64)
						if err != nil {
							errStr := fmt.Sprintln("failed to parse height:", err.Error())
							util.AddLog(errStr, comic.Title)
						} else {
							height = float64(h) * 0.0104166667
						}
					}

					panel := util.Panel{
						URL:       url,
						Image:     data,
						ImageType: imageType,
						Width:     width,
						Height:    height,
					}

					panels = append(panels, panel)
				},
			)

			// TODO: create multiple output option
			// create episode pdf
			title := episodeMap[g.Opt.StartURLs[0]]
			err := util.CreatePDF(comic.Title, title, panels)
			if err != nil {
				errStr := fmt.Sprintf("failed to create pdf: %s\n", err.Error())
				util.AddLog(errStr, comic.Title)
			}
		},
		LogDisabled: !op.Verbose,
	}).Start()
}
