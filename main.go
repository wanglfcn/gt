package main

import (
	"github.com/nsf/termbox-go"
	"fmt"
	"time"
)

var serverList *ServerList

func main() {

	title 	 := ""
	host 	 := ""
	username := ""
	password := ""

	last_g := time.Now().UnixNano()

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
					title, username, password, host = serverList.select_node()
					break mainloop

				case termbox.KeyEsc:
					break mainloop

				case termbox.KeyCtrlD:
					serverList.page_down()

				case termbox.KeyCtrlU:
					serverList.page_up()

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

					case 'g':
						if time.Now().UnixNano() - last_g < 500 * int64(time.Millisecond) {
							serverList.go_first()
						} else {
							last_g = time.Now().UnixNano()
						}

					case 'G':
						serverList.go_last()

					}
				}
			case termbox.EventResize:
				serverList.redraw()

			}
		}

	}

	termbox.Close()

	if len(host) > 0 {
		fmt.Printf("Login:\t%s (%s)\n%s\n", title, host, username + "@" + host)
		sshpass(username + "@" + host, password)
	}
}

