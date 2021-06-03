package ui

import (
	"log"
	"sync"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/comics/ui/key"
	"github.com/zerodoctor/comics/ui/view"
)

func layout(g *gocui.Gui) error {
	var err error
	maxX, maxY := g.Size()

	header := view.NewHeader(g, maxX, maxY)
	err = header.Layout()
	if err != nil {
		return err
	}

	screen := view.NewScreen(g, maxX, maxY)
	err = screen.Layout()
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
	defer g.Close()

	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan

	g.SetManagerFunc(layout)

	var wg sync.WaitGroup
	init := [...]string{"tree", "command", "screen", "header"}
	view.New(init)
	key.SetBindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}

	wg.Wait()
}
