package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

func (node *Node) Prepare() {
	node.AddRouter("/blackTest", blackTest)
	node.AddBlack("/blackTest")
	node.AddRouter("/login", login)
	node.AddRouter("/funcTest", funcTest)
	node.AddRouter("/wsTest", wsTest)
	node.AddRouter("/echoAll", echoAll)
	log.Println("prepare ready!")
}

func (node *Node) Start(port string) {
	server := &http.Server{
		Addr:              port,
		Handler:           node,
		ReadTimeout:       20 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
		//MaxHeaderBytes:    0,
	}
	// 使用自定义 handler
	//err := server.ListenAndServeTLS("server.crt", "server.key")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("build server err: ", err)
	}
}

func blackTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "error: website in blacklist!")
}

func login(w http.ResponseWriter, r *http.Request) {
	claims := &customClaims{
		"Sam",
		2,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(60*60*24) * time.Second).Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("hello"))
	if err != nil {
		log.Fatal("can not generate tokenStr! ", err)
	}
	log.Println("tokenstr generated!")
	fmt.Fprintf(w, tokenStr)
}

func funcTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("http method is not supported!")
		return
	}
	tokenStr := r.Header.Get("Authorization")
	//tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IlNhbSIsImxldmVsIjoyLCJleHAiOjE2NDY0NjIzNzMsImlzcyI6InRlc3QifQ.iLbDLT8sPjl7QpaZ2_lNzg7y_5buHR5Nzs3F5DJ_wVE"
	claims, err := ParseToken(tokenStr)
	if err != nil {
		log.Println(err)
	}
	level := claims.Level
	if level < 2 {
		fmt.Fprintf(w, "not enough authority!")
	} else {
		fmt.Fprintf(w, "funcTest success!")
	}
}
