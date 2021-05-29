package channel

import "github.com/awesome-gocui/gocui"

var (
	Shutdown = make(chan bool, 1)

	InTreeChan    = make(chan TreeData, 4)
	InHeaderChan  = make(chan HeaderData, 4)
	InScreenChan  = make(chan ScreenData, 4)
	InCommandChan = make(chan CommandData, 4)
)

type TreeData struct {
}

type HeaderType uint8

const (
	MSG HeaderType = iota
	CLOCK
)

type HeaderData struct {
	Type   HeaderType // const enum soon
	Msg    string
	Append bool
	View   *gocui.View
}

type ScreenData struct {
}

type CommandData struct {
}
