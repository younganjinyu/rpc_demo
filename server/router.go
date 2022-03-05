package server

import (
	"log"
	"net/http"
	"sync"
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

func (node *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	node.mu.RLock()
	defer node.mu.RUnlock()
	if node.Blacklist[path] == true {
		//fmt.Fprintf(w, "403 Forbidden")
		log.Println("403 Forbidden in blacklist!")
		return
	}
	// 这里取到了handler
	if handler, isOk := node.routes[path]; isOk {
		//fmt.Fprintf(w, "success!")
		handler(w, r)
		return
	}
	//fmt.Println(w, "404 Not Found")
	log.Println("404 Not Found, no path in routermap")
	return
}
