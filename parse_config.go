package main
import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
	"os"
	"io/ioutil"
)

type ServerInfo struct {
	Name 		string			`yaml:"name"`
	Ip 			string			`yaml:"ip"`
	User 		string			`yaml:"user"`
	Passwd 		string			`yaml:"password"`
	Children 	[]ServerInfo	`yaml:"children"`
}

type GtConfig struct {
	DefaultUser 	string			`yaml:"defaultUser"`
	DefaultPasswd 	string			`yaml:"defaultPassword"`
	Services 		[]ServerInfo	`yaml:"services"`
}

type Server struct {
	Name	string
	Ip	string
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

	config := GtConfig{}
	err = yaml.Unmarshal(config_content, &config)
	machines := new(Servers)

	if err != nil {
		fmt.Println("parse config encounter error: %s", err)
	}

	machines.parse_config(config.Services, config.DefaultUser, config.DefaultPasswd, 0, 0)
	machines.UpdateLines()

	return machines
}

func (this *Servers)parse_config(config []ServerInfo, defaultUser string, defaultPassword string, level int, index int) (num int) {
	if config == nil {
		return
	}

	for _, machine := range config {

		user := machine.User
		if user == "" {
			user = defaultUser
		}

		password := machine.Passwd
		if password == "" {
			password = defaultPassword
		}

		this.services = append(this.services, Server{Name: machine.Name, User: user, Ip: machine.Ip, Passwd: password, Level: level, Visible: false, Index: index, Leaf: true})
		index ++

		if machine.Children != nil && len(machine.Children) > 0 {
			this.services[index - 1].Leaf = false
			index = this.parse_config(machine.Children, defaultUser, defaultPassword, level + 1, index)
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

	for i := index + 1; i < len(this.services) && this.services[i].Level >= level; i ++ {
		if this.services[i].Level == level {
			this.services[i].Visible = true
		}
	}

	this.UpdateLines()

}

func (this *Servers) CloseNode(index int) {
	if index >= len(this.services) {
		return
	}

	level := this.services[index].Level + 1

	for i := index + 1; i < len(this.services) && this.services[i].Level >= level; i ++ {
		this.services[i].Visible = false
	}

	this.UpdateLines()

}

func (this *Servers) UpdateLines() {
	var consider, i , count int
	count = len(this.services)

	this.title_len = 0

	for index, server := range this.services {

		if len(server.Name) > this.title_len {
			this.title_len = len(server.Name)
		}

		lines := make([]string, server.Level + 1)

		for i := 0; i < server.Level + 1; i ++ {
			lines[i] = "   "
		}

		var open_status int

		open_status = '+'

		if index < len(this.services) - 1 && this.services[index + 1].Visible {
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

			if (cur_node.Visible || cur_node.Level == 0) && cur_node.Level <= consider {

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
