package key

import (
	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/comics/ui/view"
)

func SetBindings(g *gocui.Gui) {
	// global bindings
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, view.Quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, view.Quit); err != nil {
		panic(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, view.NextView); err != nil {
		panic(err)
	}

	// screen bindings
	if err := g.SetKeybinding("screen", rune('j'), gocui.ModNone, downScreen); err != nil {
		panic(err)
	}

	if err := g.SetKeybinding("screen", rune('k'), gocui.ModNone, upScreen); err != nil {
		panic(err)
	}
}
