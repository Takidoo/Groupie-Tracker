package main

import (
	"html/template"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/search.html")
	tmpl.Execute(w, nil)
}

func searchFunc(input string) {

}

func main() {
	http.HandleFunc("/", start)

	http.ListenAndServe(":80", nil)
}
