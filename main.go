package main

import (
	"GroupieTrackerModule/GroupieModule"
	"net/http"
)

func main() {
	if !GroupieModule.RsrcInit() {
		print("\nImpossible de démarrer le serveur web")
	}
	http.Handle("/ressources/", http.StripPrefix("/ressources/", http.FileServer(http.Dir("./ressources"))))
	http.HandleFunc("/", GroupieModule.SearchPage)
	http.HandleFunc("/info", GroupieModule.InfoPage)
	http.ListenAndServe(":80", nil)
}
