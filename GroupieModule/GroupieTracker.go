package GroupieModule

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

var ArtistBody *http.Response
var Err error
var Infos []GroupInfos

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data := PageData{
			Groups: SearchFunc(r.FormValue("userInput"), Infos),
		}
		tmpl, _ := template.ParseFiles("templates/search.html")
		tmpl.Execute(w, data)
	} else {
		data := PageData{
			Groups: Infos,
		}
		tmpl, _ := template.ParseFiles("templates/search.html")
		tmpl.Execute(w, data)
	}
}

func SearchFunc(input string, infos []GroupInfos) []GroupInfos {
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		if strings.Contains(infos[i].Name, input) {
			output = append(output, infos[i])
		}
	}
	return output
}
func InfoPage(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(ArtistBody.Body).Decode(&Infos)
	data := PageData{
		SearchFunc(r.URL.Query().Get("id"), Infos),
	}
	tmpl, _ := template.ParseFiles("templates/groupinfo.html")
	tmpl.Execute(w, data)
}

func RsrcInit() bool {
	ArtistBody, Err = http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if Err != nil {
		print("Impossible de charger l'API 'Artiste'")
		return false
	}
	json.NewDecoder(ArtistBody.Body).Decode(&Infos)

	return true
}
