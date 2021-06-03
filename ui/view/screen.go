package view

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/comics/ui/channel"
)

type Screen struct {
	view       *gocui.View
	g          *gocui.Gui
	maxX, maxY int
}

func NewScreen(g *gocui.Gui, maxX, maxY int) *Screen {
	return &Screen{
		g:    g,
		maxX: maxX,
		maxY: maxY,
	}
}

func (s *Screen) Layout() error {
	if v, err := s.g.SetView("screen", 0, (s.maxY/15)+1, s.maxX-1, s.maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "screen"
		v.Wrap = false
		s.view = v

		go s.PrintView()
	}

	return nil
}

func (s *Screen) Display(msg string) {
	s.g.UpdateAsync(func(g *gocui.Gui) error {
		s.view.Clear()
		fmt.Fprint(s.view, msg)
		return nil
	})
}

func (s *Screen) PrintView() {

	for data := range channel.InScreenChan {
		s.Display(data.Msg)
	}
}
