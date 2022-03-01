package serverRpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"rpcApi/common"
)

type Function struct {
}

// Set permissions based on function name

// PCal P==Public level
func (f *Function) PCal(request *common.Args, reply *int) error {
	*reply = request.Param1 * request.Param2
	return nil
}

// VCal V==VIP level
func (f *Function) VCal(request *common.Args, reply *int) error {
	if request.Level < 1 {
		return errors.New("not enough authority")
	} else {
		*reply = request.Param1 + request.Param2
		return nil
	}
}

// KCal K==King level
func (f *Function) KCal(request *common.Args, reply *int) error {
	if request.Level < 2 {
		return errors.New("not enough authority")
	} else {
		*reply = request.Param1*request.Param2 + 10
		return nil
	}
}

//func Getmanager() *common.ClientsNotifier {
//	return Manager
//}
//
//var Manager = common.NewClientsNotifier()

func SetupJsonRpcServer() {
	//Manager := common.NewClientsNotifier()
	fun := new(Function)
	err := rpc.Register(fun)
	if err != nil {
		fmt.Println("register error")
	}
	//server := rpc.NewServer()
	//err := server.RegisterName("Auservice", fun)
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print("rpc.Serve: accept:", err)
			return
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
