package main

import (
	"github.com/nsf/termbox-go"
	"fmt"
)

var serverList *ServerList

func main() {

	host := ""
	username := ""
	password := ""

	err := termbox.Init()

	if err != nil {
		panic(err)
	}


	serverList = NewServerList()

	serverList.redraw()

	mainloop:
	for {
		event := termbox.PollEvent()
		if serverList.mode == Search {
			switch event.Key {
			case termbox.KeyEnter:
				serverList.search()
				serverList.set_normal_mode()

			case termbox.KeyEsc:
				serverList.set_normal_mode()
				serverList.clear_search()

			case termbox.KeyBackspace: fallthrough
			case termbox.KeyBackspace2:
				serverList.delete_search_str()

			case termbox.KeySpace:
				serverList.add_search_str(" ")

			default:
				serverList.add_search_str(string(event.Ch))

			}

			serverList.redraw()

		} else {
			switch event.Type {
			case termbox.EventKey :
				switch event.Key {
				case termbox.KeyEnter:
					username, password, host = serverList.select_node()
					break mainloop

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

					case 'n': serverList.go_next(true)

					case 'N': serverList.go_next(false)

					case '/':
						serverList.set_search_mode()
					}
				}
			case termbox.EventResize:
				serverList.redraw()

			}
		}

	}

	termbox.Close()

	if len(host) > 0 {
		fmt.Printf("Host:%s\n", host)
		sshpass(username + "@" + host, password)
	}
}

