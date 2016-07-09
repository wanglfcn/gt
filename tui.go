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

type Mode uint

const (
	Normal Mode = iota
	Search
)

type SearchNode struct {
	search_str	string
	result 		[]int
	curr_pos	int
}

type ServerList struct {
	width 		int
	high		int
	offset 		int
	currentIndex 	int
	currentID	int
	servers		*Servers
	titles		[]Title
	mode		Mode
	searchNode 	SearchNode
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
	serverList.mode = Normal
	serverList.searchNode.search_str = ""
	serverList.searchNode.curr_pos = 0
	return serverList
}

func (this *ServerList)moveUp() {
	if this.currentIndex > 0 {
		this.currentIndex -= 1
		this.currentID = this.titles[this.currentIndex].id
		if this.currentIndex < this.offset {
			this.offset = this.currentIndex
		}
		this.redraw()
	}
}

func (this *ServerList)moveDown() {
	if this.currentIndex < (len(this.titles) - 1) {
		this.currentIndex += 1
		this.currentID = this.titles[this.currentIndex].id

		if this.currentIndex - this.offset > this.high - 8 {
			this.offset ++
		}

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

func (this *ServerList)drawText(x, y int, text string, fg, bg termbox.Attribute) {
	if x < 1 || x >= this.width - 1 || y < 1 || y >= this.high - 1 {
		return
	}

	for _, c := range text {
		if x >= this.width - 1 {
			break
		}
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
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

	this.drawText(start_pos, 2, title, fg, bg)

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

	termbox.SetCell(this.width - 9, this.high - 3, '┬', fg, bg)
	termbox.SetCell(this.width - 9, this.high - 2, '│', fg, bg)
	termbox.SetCell(this.width - 9, this.high - 1, '┴', fg, bg)

	status := "Normal";
	if (this.mode == Search) {
		status = "Search"
		fg, bg = bg, fg
	}

	this.drawText(2, this.high - 2, this.searchNode.search_str, fg, bg)

	this.drawText(this.width - 8, this.high - 2, status, fg, bg)

	return true
}

func (this *ServerList)search() {
	this.searchNode.result = this.searchNode.result[:0]

	count := len(this.servers.services)

	for i := 0; i < count; i ++ {
		server := this.servers.services[(this.currentID + i) % count]

		if (
		strings.Contains(server.Ip, this.searchNode.search_str) ||
		strings.Contains(strings.ToLower(server.Name), strings.ToLower(this.searchNode.search_str))) {
			this.searchNode.result = append(this.searchNode.result, server.Index)
		}
	}

	this.searchNode.curr_pos = -1
	this.go_next(true)
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

	for i := this.offset; i < this.high + this.offset - 7 && i < len(this.titles); i ++ {
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

func (this *ServerList)clear_search() {
	this.searchNode.search_str = ""
	this.searchNode.result = this.searchNode.result[:0]
}

func (this *ServerList)add_search_str(str string) {
	this.searchNode.search_str += str
}

func (this *ServerList)delete_search_str() {
	last := len(this.searchNode.search_str)
	if last > 0 {
		this.searchNode.search_str = this.searchNode.search_str[:last - 1]
	}
}

func (this *ServerList)set_normal_mode() {
	this.mode = Normal
}

func (this *ServerList)set_search_mode() {
	this.mode = Search
	this.clear_search()
	this.redraw()
}

func (this *SearchNode)get_index(down bool)(index int, ok bool) {
	count := len(this.result)
	if (count == 0) {
		index = -1
		ok = false
		return
	}

	ok = true
	if down {
		this.curr_pos = (this.curr_pos + 1) % count
	} else {
		this.curr_pos = (this.curr_pos - 1 + count) % count
	}

	index = this.result[this.curr_pos]
	return
}

func (this *ServerList)go_next(down bool) {
	index, ok := this.searchNode.get_index(down)

	if ok == false {
		return
	}

	//待优化

	level := this.servers.services[index].Level

	for i := index; i >= 0 && level >= 0; i -- {
		if (this.servers.services[i].Level < level) {
			this.currentID = index
			this.expandNode()
			level --
		}
	}

	for i, server := range this.titles {
		if server.id == index {
			this.currentIndex = i

			if this.currentIndex - this.offset > this.high - 8 {
				this.offset ++
			} else if this.currentIndex < this.offset {
				this.offset = this.currentIndex
			}

			this.redraw()
			return
		}

	}
}