package main

import (
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"little-computer-3/machine"
)

func menuView(gui *gocui.Gui, maxX, maxY int) error {
	if view, err := gui.SetView("menu", -1, maxY-5, maxX, maxY+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	
		_, err := gui.SetCurrentView("menu")

		if err != nil {
			return err
		}

		view.Autoscroll = false
		view.Editable = false
		view.Wrap = false
		view.Frame = false

		view.FgColor = gocui.Attribute(15 + 1)
		view.BgColor = gocui.ColorDefault

		go func() {
			for range time.Tick(time.Millisecond * 100) {
				updateMenuView(gui)
			}
		}()

		updateMenuView(gui)
	}

	return nil
}

func updateMenuView(gui *gocui.Gui) {
	return
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	vm := machine.NewLC3VM()
	if vm != nil {
		panic("filler")
	}
}

func layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()
	
	if err := menuView(gui, maxX, maxY); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}