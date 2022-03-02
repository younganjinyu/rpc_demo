package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type application struct {
	routes map[string]func(http.ResponseWriter, *http.Request)
	while  map[string]bool
}

func Create() *application {
	return &application{
		routes: make(map[string]func(http.ResponseWriter, *http.Request)),
		while:  make(map[string]bool),
	}
}

func (app *application) Router(path string, controller func(http.ResponseWriter, *http.Request)) {
	app.routes[path] = controller
}

func (app *application) Jump(path string) {
	app.while[path] = true
}

func (app *application) Start(bindPort string) {
	// 使用自定义 handler
	err := http.ListenAndServe(bindPort, app)
	if err != nil {
		log.Fatal("服务器创建失败")
	}
}

func PermissionCheck(w http.ResponseWriter, r *http.Request) (string, error) {
	// 我是公共验证权限的，每个不在白名单中的请求都会访问到我
	return "不允许通过", errors.New("不允许通过")
}

func NotFind(w http.ResponseWriter, r *http.Request) string {
	return "未找到访问的内容"
}

func Find(w http.ResponseWriter, r *http.Request) string {
	return "找到合适方法"
}

func (app *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if app.while[path] != true {
		message, err := PermissionCheck(w, r)
		if err != nil {
			w.Write([]byte(message))
			return
		}
	}
	if _, function := app.routes[path]; function {
		// 返回数据给客户端
		w.Write([]byte(Find(w, r)))
		return
	}
	// 404 未找到用户访问的地址
	w.Write([]byte(NotFind(w, r)))
	return
}

func main() {
	app := Create()
	app.Router("/black", Black)
	app.Router("/white", White)
	app.Jump("/white")
	app.Start(":8080")
}

func White(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "我是白名单")
}

func Black(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "我是黑名单")
}
