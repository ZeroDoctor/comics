package ui

import (
	"log"
	"sync"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/comics/ui/key"
	"github.com/zerodoctor/comics/ui/view"
)

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	err := view.SetHeaderView(g, maxX, maxY)
	if err != nil {
		return err
	}

	return nil
}

func Start() {

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Fatal(err)
	}

	g.Mouse = false
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan

	g.SetManagerFunc(layout)

	var wg sync.WaitGroup
	init := [...]string{"tree", "command", "screen", "header"}
	view.New(init)
	view.SetupViews(g, &wg)
	key.SetBindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}

	wg.Wait()
}
