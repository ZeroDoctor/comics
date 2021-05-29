package view

import (
	"fmt"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/comics/ui/channel"
)

type Header struct {
	view       *gocui.View
	g          *gocui.Gui
	maxX, maxY int
}

func NewHeader(g *gocui.Gui, maxX, maxY int) *Header {
	return &Header{
		g:    g,
		maxX: maxX,
		maxY: maxY,
	}
}

func (h *Header) Layout() error {
	if v, err := h.g.SetView("header", 0, 0, h.maxX-1, (h.maxY / 15), 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "header"
		v.Wrap = false
		h.view = v

		go h.PrintView()
	}

	return nil
}

func (h *Header) PrintView() {
	go h.showClock()

	for data := range channel.InHeaderChan {
		timeNow := time.Now().Format("02/01/2006 15:04:05")
		msg := ""
		clock := "Date: " + timeNow

		switch data.Type {
		case channel.MSG:
			if data.Append {
				msg = " " + data.Msg
			} else {
				msg += " " + data.Msg
			}
		}
		h.Display(clock, msg)
	}
}

func (h *Header) Display(clock, msg string) {
	h.g.UpdateAsync(func(g *gocui.Gui) error {
		h.view.Clear()
		fmt.Fprint(h.view, clock+msg)
		return nil
	})
}

func (h *Header) showClock() {
	for {
		select {
		case <-channel.Shutdown:
			return
		default:
			if h.view != nil {
				channel.InHeaderChan <- channel.HeaderData{Type: channel.CLOCK}
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
}
