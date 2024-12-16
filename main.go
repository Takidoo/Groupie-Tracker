package main

import (
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	println("gvdfgidfg")
}

func main() {
	http.HandleFunc("/", start)
	http.ListenAndServe(":80", nil)
}
