package view

import "github.com/awesome-gocui/gocui"

type Tree struct {
	view       *gocui.View
	g          *gocui.Gui
	maxX, maxY int
}

func NewTree(g *gocui.Gui, maxX, maxY int) *Tree {
	return &Tree{
		g:    g,
		maxX: maxX,
		maxY: maxY,
	}
}

func (t *Tree) Layout() error {
	if v, err := t.g.SetView("tree", 0, (t.maxY/15)+1, (t.maxX / 6), t.maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "tree"
		v.Wrap = false
		t.view = v
	}

	return nil
}
