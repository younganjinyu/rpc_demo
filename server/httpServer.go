package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	mu        sync.RWMutex
	routes    map[string]func(http.ResponseWriter, *http.Request)
	Blacklist map[string]bool
}

func Create() *Node {
	return &Node{
		routes:    make(map[string]func(http.ResponseWriter, *http.Request)),
		Blacklist: make(map[string]bool),
	}
}

func (node *Node) AddRouter(path string, controller func(http.ResponseWriter, *http.Request)) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if path == "" {
		panic("http: invalid pattern")
	}
	if controller == nil {
		panic("http: nil handler")
	}
	if _, isExist := node.routes[path]; isExist {
		panic("http: multiple registrations for " + path)
	}
	node.routes[path] = controller
}

func (node *Node) DelRouter(path string) {
	node.mu.Lock()
	defer node.mu.Unlock()
	delete(node.routes, path)
}

func (node *Node) AddBlack(path string) {
	node.mu.Lock()
	defer node.mu.Unlock()
	if path == "" {
		panic("http: invalid pattern")
	}
	if _, isExist := node.Blacklist[path]; isExist {
		panic("http: multiple registrations for " + path)
	}
	node.Blacklist[path] = true
}

func (node *Node) DelBlack(path string) {
	node.mu.Lock()
	defer node.mu.Unlock()
	delete(node.Blacklist, path)
}

func (node *Node) Start(port string) {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           node,
		ReadTimeout:       20 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
		//MaxHeaderBytes:    0,
	}
	// 使用自定义 handler
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("服务器创建失败")
	}
}

func (node *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	node.mu.RLock()
	defer node.mu.RUnlock()
	if node.Blacklist[path] == true {
		fmt.Fprintf(w, "403 Forbidden")
		return
	}
	// 这里取到了handler
	if handler, isOk := node.routes[path]; isOk {
		fmt.Fprintf(w, "success!")
		handler(w, r)
		return
	}
	fmt.Println(w, "404 Not Found")
	return
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 3 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}

func (node *Node) Prepare() {
	node.AddRouter("/blackTest", blackTest)
	node.AddBlack("/blackTest")
	node.AddRouter("/normal", normal)
	node.AddRouter("/wsTest", wsTest)
	log.Println("prepare ready!")
}

func wsTest(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Connected!"+"client's address: ", conn.RemoteAddr())
	err = conn.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println("hi error")
		return
	}
	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func blackTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "error: website in blacklist!")
}

const LEVEL2 = 2

func normal(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		levelStr := r.FormValue("level")
		level, _ := strconv.Atoi(levelStr)
		if level < LEVEL2 {
			fmt.Fprintf(w, "not permitted!")
		} else {
			fmt.Fprintf(w, "this is normal test")
		}
	} else {
		fmt.Fprintf(w, "http method is not permitted")
	}
}
