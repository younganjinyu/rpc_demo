package clientRpc

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"rpcApi/common"
)

//func GetBalance(addr string) (*big.Int, error) {
//	return big.NewInt(0), nil
//}

func SetupNewClient() {
	client, err := jsonrpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("dial error:", err)
	}
	////Manager
	//serverRpc.Getmanager().Register(client)
	//defer serverRpc.Getmanager().Deregister(client)
	// get level from DB
	args := &common.Args{Level: 1, Param1: 1, Param2: 2}
	var reply int
	err = client.Call("Function.PCal", args, &reply)
	if err != nil {
		log.Fatal("Client call error:", err)
	}
	//else {
	//	time.Sleep(60 * time.Second)
	//}
	fmt.Println(reply)
}
