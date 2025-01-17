package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

type PageData struct {
	Groups []GroupInfos
}
type GroupInfos struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

func addMj(tab *[]GroupInfos) {
	*tab = append(*tab, GroupInfos{
		ID:           53,
		Image:        "https://media.vanityfair.fr/photos/63871e8919ae25c0a47e6089/16:9/w_320%2Cc_limit/3.jpg",
		Name:         "Michael Jackson",
		Members:      []string{"Michael Jackson"},
		CreationDate: 1958,
		FirstAlbum:   "Off The Wall (les autres ça ne compte pas)",
		Locations:    "",
		ConcertDates: "",
		Relations:    "",
	})
}

func start(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	var infos []GroupInfos
	json.NewDecoder(resp.Body).Decode(&infos)
	addMj(&infos)
	if r.Method == http.MethodPost {
		infos = searchFunc(r.FormValue("userInput"), infos)
	}
	gg := PageData{
		Groups: infos,
	}
	tmpl, _ := template.ParseFiles("templates/search.html")
	tmpl.Execute(w, gg)
}

func searchFunc(input string, infos []GroupInfos) []GroupInfos {
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		if strings.Contains(infos[i].Name, input) {
			output = append(output, infos[i])
		}
	}
	return output
}
func infoPage(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	var infos []GroupInfos
	json.NewDecoder(resp.Body).Decode(&infos)
	if r.URL.Query().Get("id") == "Michael Jackson" {
		addMj(&infos)
	}
	gg := PageData{
		searchFunc(r.URL.Query().Get("id"), infos),
	}
	tmpl, _ := template.ParseFiles("templates/groupinfo.html")
	tmpl.Execute(w, gg)
}

func main() {
	http.HandleFunc("/", start)
	http.HandleFunc("/info", infoPage)
	http.ListenAndServe(":80", nil)
}
