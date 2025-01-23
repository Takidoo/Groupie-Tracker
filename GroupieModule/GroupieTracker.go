package GroupieModule

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type SearchPageData struct {
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
var ArtistErr error
var Infos []GroupInfos

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data := SearchPageData{
			Groups: SearchFunc(r.FormValue("userInput"), Infos),
		}
		tmpl, _ := template.ParseFiles("templates/search.html")
		tmpl.Execute(w, data)
	} else {
		data := SearchPageData{
			Groups: Infos,
		}
		tmpl, _ := template.ParseFiles("templates/search.html")
		tmpl.Execute(w, data)
	}
}

func SearchFunc(input string, infos []GroupInfos) []GroupInfos {
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		if strings.Contains(strings.ToLower(infos[i].Name), strings.ToLower(input)) {
			output = append(output, infos[i])
		}
	}
	return output
}
func InfoPage(w http.ResponseWriter, r *http.Request) {
	var groupId, _ = strconv.Atoi(r.URL.Query().Get("id"))
	groupId += -1
	data := GroupInfos{
		ID:           Infos[groupId].ID,
		Image:        Infos[groupId].Image,
		Name:         Infos[groupId].Name,
		Members:      Infos[groupId].Members,
		CreationDate: Infos[groupId].CreationDate,
		FirstAlbum:   Infos[groupId].FirstAlbum,
		Locations:    Infos[groupId].Locations,
		ConcertDates: Infos[groupId].ConcertDates,
		Relations:    Infos[groupId].Relations,
	}
	tmpl, _ := template.ParseFiles("templates/groupinfo.html")
	tmpl.Execute(w, data)
}

func RsrcInit() bool {
	ArtistBody, ArtistErr = http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if ArtistErr != nil {
		print("Impossible de charger l'API 'Artiste'")
		return false
	}
	json.NewDecoder(ArtistBody.Body).Decode(&Infos)

	return true
}
