package key

import "github.com/awesome-gocui/gocui"

func upScreen(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-5); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-5); err != nil {
				return err
			}
		}
	}

	return nil
}

// downScreen :
func downScreen(g *gocui.Gui, v *gocui.View) error {

	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+5); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+5); err != nil {
				return err
			}
		}
	}

	return nil
}
