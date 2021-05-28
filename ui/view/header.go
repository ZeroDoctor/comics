package view

import (
	"fmt"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/comics/ui/channel"
)

type Header struct {
	view *gocui.View
	g    *gocui.Gui

	msg   string
	clock string
}

func SetHeaderView(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("header", 0, 0, maxX-1, (maxY / 15), 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "header"
		v.Wrap = false
	}

	return nil
}

func (h *Header) PrintView(wg *sync.WaitGroup) {
	defer wg.Done()
	go h.showClock()

	for data := range channel.InHeaderChan {
		if h.view != nil {
			timeNow := time.Now().Format("02/01/2006 15:04:05")
			h.clock = "Date: " + timeNow

			switch data.Type {
			case "msg":
				if data.Boolean {
					h.msg = " " + data.String
				} else {
					h.msg += " " + data.String
				}
			}
			h.Display()
		} else if data.Type == "view" {
			h.view = data.Object.(*gocui.View)
		}
	}
}

// Display :
func (h *Header) Display() {
	h.g.UpdateAsync(func(g *gocui.Gui) error {
		h.view.Clear()
		fmt.Fprint(h.view, h.clock+h.msg)
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
				channel.InHeaderChan <- channel.Data{Type: "clock"}
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
}
