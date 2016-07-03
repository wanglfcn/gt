package main
import (
	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
)

func draw_bondary(w, h int, fg, bg termbox.Attribute, title string) bool {

	if w < 5 || h < 9 {
		return false
	}

	title_width := runewidth.StringWidth(title)

	start_pos := (w - 4 - title_width) / 2 + 2

	if start_pos < 2 {
		start_pos = 2
	}

	for _, c := range title {
		if start_pos >= w - 1 {
			break
		}
		termbox.SetCell(start_pos, 2, c, fg, bg)
		start_pos += runewidth.RuneWidth(c)
	}

	for x := 1; x < w - 1; x ++ {
		termbox.SetCell(x, 1, '─', fg, bg)
		termbox.SetCell(x, 3, '─', fg, bg)
		termbox.SetCell(x, h - 3, '─', fg, bg)
		termbox.SetCell(x, h - 1, '─', fg, bg)
	}

	for y := 1; y < h - 1; y ++ {
		termbox.SetCell(1, y, '│', fg, bg)
		termbox.SetCell(w - 2, y, '│', fg, bg)
	}

	termbox.SetCell(1, 1, '┌', fg, bg)
	termbox.SetCell(1, 3, '├', fg, bg)
	termbox.SetCell(1, h - 3, '├', fg, bg)
	termbox.SetCell(1, h - 1, '└', fg, bg)

	termbox.SetCell(w - 2, 1, '┐', fg, bg)
	termbox.SetCell(w - 2, 3, '┤', fg, bg)
	termbox.SetCell(w - 2, h - 3, '┤', fg, bg)
	termbox.SetCell(w - 2, h - 1, '┘', fg, bg)

	return true
}

var serverList *ServerList

func main() {

	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	width, high := termbox.Size()
	serverList = NewServerList(width, high)

	draw_bondary(width, high, termbox.ColorWhite, termbox.ColorDefault, "Machine list")
	serverList.redraw()
	termbox.Flush()

	mainloop:
	for {
		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyEsc:
				break mainloop
			default:
				width, high = termbox.Size()
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				draw_bondary(width, high, termbox.ColorWhite, termbox.ColorDefault, "machine list")
			}

			switch event.Ch {
			case 'Q': fallthrough
			case 'q': break mainloop
			case 'J': fallthrough
			case 'j':
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				serverList.moveDown()
				serverList.redraw()
				draw_bondary(width, high, termbox.ColorWhite, termbox.ColorDefault, "machine list")

			case 'K': fallthrough
			case 'k':
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				serverList.moveUp()
				serverList.redraw()
				draw_bondary(width, high, termbox.ColorWhite, termbox.ColorDefault, "machine list")

			}

		case termbox.EventResize:
			width, high = termbox.Size()
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			draw_bondary(width, high, termbox.ColorWhite, termbox.ColorDefault, "machine list")
		}

		termbox.Flush()

	}

}

