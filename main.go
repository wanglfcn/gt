package main
import (
	"github.com/nsf/termbox-go"
)

var serverList *ServerList

func main() {

	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	serverList = NewServerList()

	serverList.redraw()

	mainloop:
	for {

		switch event := termbox.PollEvent(); event.Ch {
		case 'Q': fallthrough
		case 'q': break mainloop

		case 'J': fallthrough
		case 'j': serverList.moveDown()

		case 'K': fallthrough
		case 'k': serverList.moveUp()

		case 'I': fallthrough
		case 'i': serverList.expandNode()

		case 'O': fallthrough
		case 'o': serverList.closeNode()

		}

	}

}

