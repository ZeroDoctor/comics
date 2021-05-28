package view

import (
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/comics/ui/channel"
)

var vl *List

type List struct {
	list        [4]string
	currentView int
}

func New(views [4]string) {
	vl = &List{list: views, currentView: 0}
}

func SetupViews(g *gocui.Gui, wg *sync.WaitGroup) {
	header := &Header{g: g}

	wg.Add(1)
	go header.PrintView(wg)

	go waitForViews(g)
}

func waitForViews(g *gocui.Gui) {
	for {
		time.Sleep(100 * time.Millisecond)

		header, err := g.View("header")
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		channel.InHeaderChan <- channel.Data{Type: "view", Object: header}

		return
	}
}

func SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	return g.SetViewOnTop(name)
}

func NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (vl.currentView + 1) % (len(vl.list) - 1)
	name := vl.list[nextIndex]

	if _, err := SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	vl.currentView = nextIndex
	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {

	channel.Shutdown <- true

	time.Sleep(750 * time.Millisecond)

	close(channel.InHeaderChan)
	close(channel.InScreenChan)
	close(channel.InCommandChan)

	return gocui.ErrQuit
}
