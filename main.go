package main
<<<<<<< HEAD
=======

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
>>>>>>> 13890fc6020f49eb09ed6bdc5dbe25ec271df246
