package main

import (
	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
	_ "unicode/utf8"
	"fmt"
	"strings"
)

var colors = []termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorRed,
	termbox.ColorGreen,
	termbox.ColorYellow,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorCyan,
	termbox.ColorWhite,
}

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

func NewServerList(width, hight int) *ServerList {

	serverList = new(ServerList)
	serverList.servers = NewServices()
	serverList.width = width
	serverList.high = hight
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
			x += runewidth.RuneWidth(c)
		}
	}

}

func (this *ServerList)redraw() {
	if this.currentIndex  < this.offset {
		this.offset = this.currentIndex
	}

	if this.currentIndex >= (this.offset + this.high) {
		this.offset = this.currentIndex - this.high + 1
	}

	for i := this.offset; i < this.high + this.offset && i < len(this.titles); i ++ {
		this.drawLine(i, this.offset, i == this.currentIndex)
	}
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
