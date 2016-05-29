package main
import (
	"fmt"
	"log"
	"github.com/jroimartin/gocui"
)

func main() {
	var servers *Servers = NewServices()

	gui := gocui.NewGui()

	if err := gui.Init(); err != nil {
		log.Panic(err)
	}

	defer gui.Close()

	gui.SetLayout(layout)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if  err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	fmt.Println(servers.services)

}

func layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView("hello", maxX/2 - 7, maxY/2, maxX/2 + 7, maxY/2 + 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "hello world")
	}
	return nil
}

func quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
