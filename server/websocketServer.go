package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 3 * time.Second,
	ReadBufferSize:   4096,
	WriteBufferSize:  4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsTest(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade err: ", err)
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read msg err: ", err)
			break
		}
		log.Println("msgTtpe: ", msgType, "  msg: ", msg)
		msg = []byte("received! ")
		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("write err: ", err)
			break
		}
	}
}

func echoAll(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade err: ", err)
		return
	}
	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read err: ", err)
			break
		}
		if mt == 1 {
			err = conn.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write err: ", err)
				break
			}
		}
	}
}
