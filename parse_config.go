package main
import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"os"
	"io/ioutil"
)

//parse content

type Server struct {
	Name	string
	Ip	string
	User	string
	Passwd	string
	Level	int
	Visible	bool
	index 	int
}

type Servers struct {
	services []Server
}

func main(){

	config_path := os.Getenv("GT_CONFIG")
	config_content, err := ioutil.ReadFile(config_path)
	if err != nil {
		fmt.Println("read file error %s\n", err)
		os.Exit(1)
	}

	config, err := simplejson.NewJson(config_content)
	machines := new(Servers)

	if err != nil {
		fmt.Println("parse config encounter error: %s \n ", err)
	}

	machines.parse_config(config, 0, 0)
	fmt.Println("res: %v", machines.services)

	machines.OpenNode(1)
	fmt.Println("res: %v", machines.services)
}

func (this *Servers)parse_config(config *simplejson.Json, level int, index int) (num int) {
	var i int = 0
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

			this.services = append(this.services, Server{Name: name, User: user, Ip: ip, Passwd: passwd, Level: level, Visible: false, index: index})
			index ++
		}

		if child, exist := machine.CheckGet("child"); exist {
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
	level := this.services[index].Level

	for i := index + 1; i < len(this.services) && this.services[i].Level > level; i ++ {
		this.services[i].Visible = true
	}

}

func (this *Servers) CloseNode(index int) {
	if index >= len(this.services) {
		return
	}

	this.services[index].Visible = false
	level := this.services[index].Level

	for i := index + 1; i < len(this.services) && this.services[i].Level > level; i ++ {
		this.services[i].Visible = false
	}

}
