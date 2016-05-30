package main
import (
	"fmt"
	"log"
	"strings"
	"github.com/jroimartin/gocui"
)

var servers *Servers

func main() {
	servers = NewServices()

	gui := gocui.NewGui()

	if err := gui.Init(); err != nil {
		log.Panic(err)
	}

	defer gui.Close()
	fmt.Println("max length:", servers.title_len)

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
	if v, err := gui.SetView("hello", maxX/2 - 17, maxY/2 - 12, maxX/2 + 27, maxY/2 + 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		for _, server := range servers.services {
			blanks := strings.Repeat(" ", servers.title_len - len(server.Name))
			fmt.Fprintf(v, "%s %s%s\t\t%s\n", server.Lines, server.Name, blanks, server.Ip)
		}
	}
	return nil
}

func quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
