package main

import (
	"github.com/nsf/termbox-go"
	"fmt"
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
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey :
			switch event.Key {
			case termbox.KeyEnter:
				username, password, ip := serverList.select_node()
				fmt.Println("username: %s\tpassword=%sip=%s", username, password, ip)
				//break mainloop
			case termbox.KeyEsc:
				break mainloop

			default:
				switch event.Ch {
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

	}

}

