package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func wake(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
