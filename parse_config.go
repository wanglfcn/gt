package main
import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"strings"
	"os"
	"io/ioutil"
)


type Server struct {
	Name	string
	Ip		string
	User	string
	Passwd	string
	Level	int
	Visible	bool
	Index 	int
	Lines	string
	Leaf	bool
}

type Servers struct {
	services	[]Server
	title_len	int
}

func NewServices() *Servers{

	config_path := os.Getenv("GT_CONFIG")
	config_content, err := ioutil.ReadFile(config_path)
	if err != nil {
		fmt.Println("read file error %s", err)
		os.Exit(1)
	}

	config, err := simplejson.NewJson(config_content)
	machines := new(Servers)

	if err != nil {
		fmt.Println("parse config encounter error: %s", err)
	}

	machines.parse_config(config, 0, 0)
	machines.UpdateLines()

	return machines
}

func (this *Servers)parse_config(config *simplejson.Json, level int, index int) (num int) {
	var i int = 0
	this.title_len = 0
	for true {
		machine := config.GetIndex(i)
		i++

		names_json, exist := machine.CheckGet("name")
		if !exist {
			break
		}

		names := names_json.MustStringArray()
		ips := machine.Get("ip").MustStringArray()
		users := machine.Get("user").MustStringArray()
		passwds := machine.Get("passwd").MustStringArray()

		var name string
		var user string
		var passwd string

		for array_index, ip := range ips {
			if len(names) < array_index {
				name = names[0]
			} else {
				name = names[array_index]
			}

			if len(users) < array_index {
				user = users[0]
				passwd = passwds[0]
			} else {
				user = users[array_index]
				passwd = passwds[array_index]
			}

			this.services = append(this.services, Server{Name: name, User: user, Ip: ip, Passwd: passwd, Level: level, Visible: true, Index: index, Leaf: true})
			if len(name) > this.title_len {
				this.title_len = len(name)
			}
			index ++
		}

		if child, exist := machine.CheckGet("child"); exist {
			this.services[index - 1].Leaf = false
			index = this.parse_config(child, level + 1, index)
		}
	}
	return index
}

func (this *Servers) OpenNode(index int) {
	if index >= len(this.services) {
		return
	}

	this.services[index].Visible = true
	level := this.services[index].Level + 1

	for i := index + 1; i < len(this.services) && this.services[i].Level == level; i ++ {
		this.services[i].Visible = true
	}

	this.UpdateLines()

}

func (this *Servers) CloseNode(index int) {
	if index >= len(this.services) {
		return
	}

	this.services[index].Visible = false
	level := this.services[index].Level + 1

	for i := index + 1; i < len(this.services) && this.services[i].Level == level; i ++ {
		this.services[i].Visible = false
	}

	this.UpdateLines()

}

func (this *Servers) UpdateLines() {
	var consider, i , count int
	count = len(this.services)
	for index, server := range this.services {
		lines := make([]string, server.Level + 1)

		var open_status int

		open_status = '+'

		if server.Visible {
			open_status = '-'
		}

		if server.Leaf {
			open_status = 'x'
		}

		lines[server.Level] = fmt.Sprintf("[%c]", open_status)

		if server.Level == 0 {
			this.services[index].Lines = lines[0]
			continue
		}

		consider = server.Level
		lines[server.Level - 1] = " └─"

		for i = index + 1; i < count && consider > 0; i ++ {
			cur_node := this.services[i]

			if cur_node.Visible || cur_node.Level == 0 && cur_node.Level <= consider {

				if (cur_node.Level - 1 >= 0) {
					lines[cur_node.Level - 1] = " │ "
				}

				if cur_node.Level == server.Level {
					lines[server.Level - 1] = " ├─"
				}

				consider = cur_node.Level - 1

			}

		}

		this.services[index].Lines = strings.Join(lines, "")

	}
}
