package mylist

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/zerodoctor/comics/dbhandler"
	"github.com/zerodoctor/graphql"
)

var apiAnilist *API = nil
var once sync.Once

type API struct {
	lastRequest   time.Time
	rateRemaining int
	rateLimit     int
}

func init() {
	once.Do(func() {
		apiAnilist = newApi()
	})
}

func newApi() *API {
	if apiAnilist != nil {
		return apiAnilist
	}

	rateRemaining, rateLimit, lastRequest, _ := dbhandler.ConnectLite().FetchRates()

	apiAnilist = &API{
		rateRemaining: rateRemaining,
		rateLimit:     rateLimit, // * NOTE: change this do be dynamic
		lastRequest:   lastRequest,
	}

	return apiAnilist
}

func (api API) FetchAnilist(url string, req *graphql.Request, media *QMediaAnilistResp) error {
	if time.Since(api.lastRequest) > time.Minute {
		api.lastRequest = time.Now()
	}

	if api.rateRemaining < 2 {
		time.Sleep(time.Since(api.lastRequest))
	}

	client := graphql.NewClient(url)
	ctx := context.Background()
	resp, err := client.Run(ctx, req, &media)
	if err != nil {
		log.Println(err)
		return err
	}

	api.rateLimit, _ = strconv.Atoi(resp.Header["X-Ratelimit-Limit"][0])
	api.rateRemaining, _ = strconv.Atoi(resp.Header["X-Ratelimit-Remaining"][0])

	go func() {
		err := dbhandler.ConnectLite().UpdateRates(api.rateLimit, api.rateRemaining, api.lastRequest)
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}
