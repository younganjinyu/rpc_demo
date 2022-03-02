package main

import (
	"rpcApi/server"
)

func main() {
	node := server.Create()
	node.Prepare()
	node.Start(":8080")
}
