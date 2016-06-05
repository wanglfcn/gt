package main

var serverList *ServerList

func main() {
	serverList = NewServerList()

	serverList.moveDown()

	serverList.moveDown()

	serverList.moveDown()

	serverList.expandNode()
	serverList.moveDown()

	serverList.expandNode()
}

