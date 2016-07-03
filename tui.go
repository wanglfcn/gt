package main

import (
	"github.com/nsf/termbox-go"
	"fmt"
	"strings"
	"github.com/mattn/go-runewidth"
)

type Title struct {
	id 	int
	title 	string
}

type ServerList struct {
	width 		int
	high		int
	offset 		int
	currentIndex 	int
	currentID	int
	servers		*Servers
	titles		[]Title
}

func NewServerList() *ServerList {

	serverList = new(ServerList)
	serverList.servers = NewServices()
	serverList.width = 0
	serverList.high = 0
	serverList.offset = 0
	serverList.currentIndex = 0
	serverList.currentID = 0
	serverList.updateTitles()
	serverList.redraw()
	return serverList
}

func (this *ServerList)moveUp() {
	if this.currentIndex > 0 {
		this.currentIndex -= 1
		this.currentID = this.titles[this.currentIndex].id
		this.redraw()
	}
}

func (this *ServerList)moveDown() {
	if this.currentIndex < (len(this.titles) - 1) {
		this.currentIndex += 1
		this.currentID = this.titles[this.currentIndex].id
		this.redraw()
	}
}

func (this *ServerList)expandNode() {
	this.servers.OpenNode(this.currentID)
	this.updateTitles()
	this.redraw()

}

func (this *ServerList)closeNode() {
	this.servers.CloseNode(this.currentID)
	this.updateTitles()
	this.redraw()

}

func (this *ServerList)drawLine(index int, offset int, selected bool) {
	if index < len(this.titles) {
		background := termbox.ColorDefault

		if selected {
			background = termbox.ColorGreen
			termbox.SetCell(2, index - offset + 4, '>', termbox.ColorWhite, background)
		}

		x := 3
		for _, c := range this.titles[index].title {
			termbox.SetCell(x, index - offset + 4, c, termbox.ColorWhite, background)
			x += 1
		}
	}

}

func (this *ServerList)boundary(fg, bg termbox.Attribute, title string) bool {

	if this.width < 5 || this.high < 9 {
		return false
	}

	title_width := runewidth.StringWidth(title)

	start_pos := (this.width - 4 - title_width) / 2 + 2

	if start_pos < 2 {
		start_pos = 2
	}

	for _, c := range title {
		if start_pos >= this.width - 1 {
			break
		}
		termbox.SetCell(start_pos, 2, c, fg, bg)
		start_pos += runewidth.RuneWidth(c)
	}

	for x := 1; x < this.width - 1; x ++ {
		termbox.SetCell(x, 1, '─', fg, bg)
		termbox.SetCell(x, 3, '─', fg, bg)
		termbox.SetCell(x, this.high - 3, '─', fg, bg)
		termbox.SetCell(x, this.high - 1, '─', fg, bg)
	}

	for y := 1; y < this.high - 1; y ++ {
		termbox.SetCell(1, y, '│', fg, bg)
		termbox.SetCell(this.width - 2, y, '│', fg, bg)
	}

	termbox.SetCell(1, 1, '┌', fg, bg)
	termbox.SetCell(1, 3, '├', fg, bg)
	termbox.SetCell(1, this.high - 3, '├', fg, bg)
	termbox.SetCell(1, this.high - 1, '└', fg, bg)

	termbox.SetCell(this.width - 2, 1, '┐', fg, bg)
	termbox.SetCell(this.width - 2, 3, '┤', fg, bg)
	termbox.SetCell(this.width - 2, this.high - 3, '┤', fg, bg)
	termbox.SetCell(this.width - 2, this.high - 1, '┘', fg, bg)

	return true
}


func (this *ServerList)redraw() {
	this.width, this.high = termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if this.currentIndex  < this.offset {
		this.offset = this.currentIndex
	}

	if this.currentIndex >= (this.offset + this.high) {
		this.offset = this.currentIndex - this.high + 1
	}

	for i := this.offset; i < this.high + this.offset && i < len(this.titles); i ++ {
		this.drawLine(i, this.offset, i == this.currentIndex)
	}

	this.boundary(termbox.ColorWhite, termbox.ColorDefault, "Machine list")
	termbox.Flush()
}

func (this *ServerList)updateTitles() {

	this.titles = this.titles[:0]
	for _, server := range this.servers.services {
		if server.Level == 0 || server.Visible {
			blanks := strings.Repeat(" ", this.servers.title_len - len(server.Name))
			title := fmt.Sprintf("%s %s%s\t\t%s\n", server.Lines, server.Name, blanks, server.Ip)
			this.titles = append(this.titles, Title{id: server.Index, title: title})
		}
	}
}

func (this *ServerList)select_node() (username, password, ip string) {
	username = ""
	password = ""
	ip 	 = ""
	if this.currentIndex >= 0 && this.currentIndex < len(this.servers.services) {
		username = this.servers.services[this.currentIndex].User
		password = this.servers.services[this.currentIndex].Passwd
		ip 	 = this.servers.services[this.currentIndex].Ip

		info := fmt.Sprintf("username=%s\tpassword=%s\tip=%s", username, password, ip)

		for _, c := range info {
			termbox.SetCell(10, 30, c, termbox.ColorWhite, termbox.ColorDefault)
		}
		termbox.Flush()
	}
	return
}